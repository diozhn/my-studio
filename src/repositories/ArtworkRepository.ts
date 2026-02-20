import { prisma } from "../database/prisma";
import { Artwork, Prisma } from "@prisma/client";

export interface ArtworkFilters {
  title?: string;
  author?: string;
  from?: string;
  to?: string;
  sort?: "likes_desc" | "likes_asc" | "created_asc" | "created_desc";
}

export interface PaginationParams {
  page: number;
  limit: number;
}

export class ArtworkRepository {
  async create(data: {
    title: string;
    caption: string;
    imageUrl: string;
    userId: number;
  }): Promise<Artwork> {
    return prisma.artwork.create({
      data,
    });
  }

  async findById(id: number): Promise<Artwork | null> {
    return prisma.artwork.findUnique({
      where: { id },
    });
  }

  async findAll(
    filters: ArtworkFilters = {},
    pagination?: PaginationParams,
  ): Promise<{ artworks: Artwork[]; total: number }> {
    const where: Prisma.ArtworkWhereInput = {};

    if (filters.title) {
      where.title = { contains: filters.title, mode: "insensitive" };
    }

    if (filters.author) {
      where.userId = parseInt(filters.author);
    }

    if (filters.from && filters.to) {
      where.createdAt = {
        gte: new Date(filters.from),
        lte: new Date(filters.to),
      };
    }

    const orderBy: Prisma.ArtworkOrderByWithRelationInput =
      filters.sort === "likes_desc"
        ? { likes: "desc" }
        : filters.sort === "likes_asc"
          ? { likes: "asc" }
          : filters.sort === "created_asc"
            ? { createdAt: "asc" }
            : { createdAt: "desc" };

    const skip = pagination ? (pagination.page - 1) * pagination.limit : undefined;
    const take = pagination?.limit;

    const findManyOptions: any = {
      where,
      orderBy,
    };

    if (skip !== undefined) {
      findManyOptions.skip = skip;
    }
    if (take !== undefined) {
      findManyOptions.take = take;
    }

    const [artworks, total] = await Promise.all([
      prisma.artwork.findMany(findManyOptions),
      prisma.artwork.count({ where }),
    ]);

    return { artworks, total };
  }

  async findByUserId(userId: number): Promise<Artwork[]> {
    return prisma.artwork.findMany({
      where: { userId },
      orderBy: { createdAt: "desc" },
    });
  }

  async getTopArtworks(limit: number): Promise<Artwork[]> {
    return prisma.artwork.findMany({
      orderBy: { likes: "desc" },
      take: limit,
    });
  }

  async update(
    id: number,
    data: { title?: string; caption?: string },
  ): Promise<Artwork> {
    return prisma.artwork.update({
      where: { id },
      data,
    });
  }

  async incrementLikes(id: number): Promise<Artwork> {
    return prisma.artwork.update({
      where: { id },
      data: { likes: { increment: 1 } },
    });
  }

  async delete(id: number): Promise<void> {
    await prisma.artwork.delete({
      where: { id },
    });
  }
}

