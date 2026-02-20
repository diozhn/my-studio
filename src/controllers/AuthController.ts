import { Request, Response } from "express";
import { AuthUseCase } from "../useCases/AuthUseCase";

export class AuthController {
  constructor(private authUseCase: AuthUseCase) {}

  async register(req: Request, res: Response) {
    try {
      const { username, email, password } = req.body;

      if (!username || !email || !password) {
        return res.status(400).json({ error: "Missing fields" });
      }

      const result = await this.authUseCase.register({
        username,
        email,
        password,
      });

      return res.json(result);
    } catch (error: any) {
      console.error(error);
      if (error.message === "Email already registered") {
        return res.status(400).json({ error: error.message });
      }
      return res.status(500).json({ error: "Error creating user" });
    }
  }

  async login(req: Request, res: Response) {
    try {
      const { email, password } = req.body;

      if (!email || !password) {
        return res.status(400).json({ error: "Invalid JSON" });
      }

      const result = await this.authUseCase.login(email, password);
      return res.json(result);
    } catch (error: any) {
      console.error(error);
      return res.status(401).json({ error: error.message || "User or password invalid" });
    }
  }

  async refreshToken(req: Request, res: Response) {
    try {
      const { refresh_token } = req.body;

      if (!refresh_token) {
        return res.status(400).json({ error: "Invalid JSON" });
      }

      const result = await this.authUseCase.refreshToken(refresh_token);
      return res.json(result);
    } catch (error: any) {
      console.error(error);
      return res.status(401).json({ error: error.message || "Invalid refresh token" });
    }
  }
}

