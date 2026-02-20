import { ArtworkRepository, ArtworkFilters, PaginationParams } from "../repositories/ArtworkRepository";
import fs from "fs/promises";
import path from "path";

export class ArtworkUseCase {
  constructor(private artworkRepository: ArtworkRepository) {}

  async create(data: {
    title: string;
    caption: string;
    imageUrl: string;
    userId: number;
  }) {
    return this.artworkRepository.create(data);
  }

  async getById(id: number) {
    const artwork = await this.artworkRepository.findById(id);
    if (!artwork) {
      throw new Error("Artwork not found");
    }
    return artwork;
  }

  async list(filters: ArtworkFilters = {}, pagination?: PaginationParams) {
    return this.artworkRepository.findAll(filters, pagination);
  }

  async getTopArtworks(limit: number = 10) {
    return this.artworkRepository.getTopArtworks(limit);
  }

  async getUserArtworks(userId: number) {
    return this.artworkRepository.findByUserId(userId);
  }

  async update(
    id: number,
    userId: number,
    data: { title?: string; caption?: string },
  ) {
    const artwork = await this.artworkRepository.findById(id);
    if (!artwork) {
      throw new Error("Artwork not found");
    }

    if (artwork.userId !== userId) {
      throw new Error("You are not the owner of this artwork");
    }

    return this.artworkRepository.update(id, data);
  }

  async delete(id: number, userId: number) {
    const artwork = await this.artworkRepository.findById(id);
    if (!artwork) {
      throw new Error("Artwork not found");
    }

    if (artwork.userId !== userId) {
      throw new Error("You are not the owner of this artwork");
    }

    // Deletar arquivo de imagem se existir
    if (artwork.imageUrl) {
      try {
        const filePath = path.join(process.cwd(), artwork.imageUrl);
        await fs.unlink(filePath);
      } catch (error) {
        // Ignorar erro se arquivo não existir
        console.error("Error deleting image file:", error);
      }
    }

    await this.artworkRepository.delete(id);
  }

  async like(id: number) {
    const artwork = await this.artworkRepository.findById(id);
    if (!artwork) {
      throw new Error("Artwork not found");
    }

    const updated = await this.artworkRepository.incrementLikes(id);
    return {
      message: "Artwork liked successfully",
      likes: updated.likes,
    };
  }

  async generateGalleryHTML() {
    const { artworks } = await this.artworkRepository.findAll(
      {},
      undefined,
    );

    let html = "<html><body><h1>Galeria de Artes</h1><div style='display: flex; flex-wrap: wrap;'>";

    for (const art of artworks) {
      html += `
        <div style='margin:10px; text-align:center'>
          <img src='${art.imageUrl}' style='max-width:200px; max-height:200px; display:block;' />
          <h3>${art.title}</h3>
          <p>${art.caption}</p>
          <p>Likes: ${art.likes}</p>
        </div>
      `;
    }

    html += "</div></body></html>";
    return html;
  }
}

