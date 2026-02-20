import { prisma } from "../database/prisma";
import { User } from "@prisma/client";

export class UserRepository {
  async create(data: {
    username: string;
    email: string;
    password: string;
  }): Promise<User> {
    return prisma.user.create({
      data,
    });
  }

  async findByEmail(email: string): Promise<User | null> {
    return prisma.user.findUnique({
      where: { email },
    });
  }

  async findById(id: number): Promise<User | null> {
    return prisma.user.findUnique({
      where: { id },
    });
  }

  async updateRefreshToken(
    userId: number,
    refreshToken: string | null,
  ): Promise<User> {
    return prisma.user.update({
      where: { id: userId },
      data: { refreshToken },
    });
  }

  async update(userId: number, data: { username?: string; password?: string }): Promise<User> {
    return prisma.user.update({
      where: { id: userId },
      data,
    });
  }

  async findBySocialId(
    provider: "google" | "instagram" | "twitter",
    socialId: string,
  ): Promise<User | null> {
    const where =
      provider === "google"
        ? { googleId: socialId }
        : provider === "instagram"
          ? { instagramId: socialId }
          : { twitterId: socialId };

    return prisma.user.findFirst({
      where,
    });
  }

  async createOrUpdateFromSocial(
    provider: "google" | "instagram" | "twitter",
    data: {
      username: string;
      email: string;
      socialId: string;
    },
  ): Promise<User> {
    const existing = await this.findBySocialId(provider, data.socialId);

    if (existing) {
      return existing;
    }

    const updateData =
      provider === "google"
        ? { googleId: data.socialId }
        : provider === "instagram"
          ? { instagramId: data.socialId }
          : { twitterId: data.socialId };

    return prisma.user.create({
      data: {
        username: data.username,
        email: data.email,
        password: "", // Social login não precisa de senha
        ...updateData,
      },
    });
  }
}

