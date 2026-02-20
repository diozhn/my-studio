import bcrypt from "bcrypt";
import { UserRepository } from "../repositories/UserRepository";
import { ArtworkRepository } from "../repositories/ArtworkRepository";

export class UserUseCase {
  constructor(
    private userRepository: UserRepository,
    private artworkRepository: ArtworkRepository,
  ) {}

  async getProfile(id: number) {
    const user = await this.userRepository.findById(id);
    if (!user) {
      throw new Error("User not found");
    }

    return {
      id: user.id,
      username: user.username,
      created_at: user.createdAt,
    };
  }

  async updateProfile(
    userId: number,
    targetUserId: number,
    data: { username?: string; password?: string },
  ) {
    if (userId !== targetUserId) {
      throw new Error("You can only update your own profile");
    }

    const updateData: { username?: string; password?: string } = {};

    if (data.username) {
      updateData.username = data.username;
    }

    if (data.password) {
      updateData.password = await bcrypt.hash(data.password, 14);
    }

    await this.userRepository.update(userId, updateData);

    return { message: "Profile updated successfully" };
  }

  async getUserArtworks(userId: number) {
    return this.artworkRepository.findByUserId(userId);
  }
}

