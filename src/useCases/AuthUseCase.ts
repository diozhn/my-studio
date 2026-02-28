import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";
import { UserRepository } from "../repositories/UserRepository";

const jwtSecret = (process.env.JWT_SECRET || "changeme") as string;

interface JwtPayload {
  user_id: number;
  exp: number;
}

export class AuthUseCase {
  constructor(private userRepository: UserRepository) {}

  async register(data: {
    username: string;
    email: string;
    password: string;
  }) {
    // Verificar se email já existe
    const existingUser = await this.userRepository.findByEmail(data.email);
    if (existingUser) {
      throw new Error("Email already registered");
    }

    const hashed = await bcrypt.hash(data.password, 14);
    const user = await this.userRepository.create({
      username: data.username,
      email: data.email,
      password: hashed,
    });

    return {
      message: "User creating success",
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
      },
    };
  }

  async login(email: string, password: string) {
    const user = await this.userRepository.findByEmail(email);
    if (!user) {
      throw new Error("User or password invalid");
    }

    const isValid = await bcrypt.compare(password, user.password);
    if (!isValid) {
      throw new Error("User or password invalid");
    }

    const token = jwt.sign(
      { user_id: user.id, exp: Math.floor(Date.now() / 1000) + 72 * 60 * 60 },
      jwtSecret,
    );

    const refreshToken = jwt.sign(
      {
        user_id: user.id,
        exp: Math.floor(Date.now() / 1000) + 7 * 24 * 60 * 60,
      },
      jwtSecret,
    );

    await this.userRepository.updateRefreshToken(user.id, refreshToken);

    return {
      is_superuser: user.superuser,
      token,
      refresh_token: refreshToken,
    };
  }

  async refreshToken(refreshToken: string) {
    try {
      const decoded = jwt.verify(refreshToken, jwtSecret) as JwtPayload;
      const user = await this.userRepository.findById(decoded.user_id);

      if (!user || user.refreshToken !== refreshToken) {
        throw new Error("Invalid refresh token");
      }

      const newToken = jwt.sign(
        {
          user_id: user.id,
          exp: Math.floor(Date.now() / 1000) + 72 * 60 * 60,
        },
        jwtSecret,
      );

      return { token: newToken };
    } catch (error) {
      throw new Error("Invalid refresh token");
    }
  }

  verifyToken(token: string): JwtPayload {
    try {
      return jwt.verify(token, jwtSecret) as JwtPayload;
    } catch (error) {
      throw new Error("Invalid token");
    }
  }
}

