import { Request, Response } from "express";
import { ArtworkUseCase } from "../useCases/ArtworkUseCase";
import multer from "multer";
import path from "path";
import fs from "fs";

// Configurar multer para upload
const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    const uploadDir = path.join(process.cwd(), "uploads");
    if (!fs.existsSync(uploadDir)) {
      fs.mkdirSync(uploadDir, { recursive: true });
    }
    cb(null, uploadDir);
  },
  filename: (req, file, cb) => {
    const uniqueName = `${Date.now()}_${file.originalname}`;
    cb(null, uniqueName);
  },
});

export const upload = multer({ storage });

export class ArtworkController {
  constructor(private artworkUseCase: ArtworkUseCase) {}

  async create(req: Request, res: Response) {
    try {
      const userId = (req as any).user_id as number;
      const { title, caption } = req.body;
      const file = req.file;

      if (!title || !file) {
        return res.status(400).json({ error: "Bad request: title and image are required" });
      }

      const imageUrl = `/uploads/${file.filename}`;

      const artwork = await this.artworkUseCase.create({
        title,
        caption: caption || "",
        imageUrl,
        userId,
      });

      return res.status(201).json(artwork);
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Failed to create artwork: " + error.message });
    }
  }

  async getById(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const id = parseInt(idParam);
      if (isNaN(id)) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const artwork = await this.artworkUseCase.getById(id);
      return res.json(artwork);
    } catch (error: any) {
      if (error.message === "Artwork not found") {
        return res.status(404).json({ error: error.message });
      }
      return res.status(500).json({ error: "Error fetching artwork" });
    }
  }

  async list(req: Request, res: Response) {
    try {
      const page = parseInt(req.query.page as string) || 1;
      const limit = parseInt(req.query.limit as string) || 10;
      const title = req.query.title as string;
      const author = req.query.author as string;
      const from = req.query.from as string;
      const to = req.query.to as string;
      const sort = req.query.sort as any;

      const filters = { title, author, from, to, sort };
      const pagination = { page, limit };

      const { artworks, total } = await this.artworkUseCase.list(filters, pagination);

      return res.json({
        page,
        limit,
        results: artworks,
        count: artworks.length,
        total,
      });
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Error fetching artworks" });
    }
  }

  async getTopArtworks(req: Request, res: Response) {
    try {
      const limit = parseInt(req.query.limit as string) || 10;
      const artworks = await this.artworkUseCase.getTopArtworks(limit);
      return res.json(artworks);
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Failed to retrieve top artworks" });
    }
  }

  async getGallery(req: Request, res: Response) {
    try {
      const html = await this.artworkUseCase.generateGalleryHTML();
      return res.type("html").send(html);
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Failed to retrieve artworks" });
    }
  }

  async update(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const id = parseInt(idParam);
      if (isNaN(id)) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const userId = (req as any).user_id as number;
      const { title, caption } = req.body;

      const artwork = await this.artworkUseCase.update(id, userId, {
        title,
        caption,
      });

      return res.json(artwork);
    } catch (error: any) {
      if (error.message === "Artwork not found") {
        return res.status(404).json({ error: error.message });
      }
      if (error.message === "You are not the owner of this artwork") {
        return res.status(403).json({ error: error.message });
      }
      return res.status(500).json({ error: "Failed to update artwork" });
    }
  }

  async delete(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const id = parseInt(idParam);
      if (isNaN(id)) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const userId = (req as any).user_id as number;

      await this.artworkUseCase.delete(id, userId);
      return res.sendStatus(204);
    } catch (error: any) {
      if (error.message === "Artwork not found") {
        return res.status(404).json({ error: error.message });
      }
      if (error.message === "You are not the owner of this artwork") {
        return res.status(403).json({ error: error.message });
      }
      return res.status(500).json({ error: "Failed to delete artwork" });
    }
  }

  async like(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const id = parseInt(idParam);
      if (isNaN(id)) {
        return res.status(400).json({ error: "Invalid artwork ID" });
      }
      const result = await this.artworkUseCase.like(id);
      return res.json(result);
    } catch (error: any) {
      if (error.message === "Artwork not found") {
        return res.status(404).json({ error: error.message });
      }
      return res.status(500).json({ error: "Failed to like artwork" });
    }
  }

  async getFiltered(req: Request, res: Response) {
    try {
      const title = req.query.title as string;
      const author = req.query.author as string;
      const from = req.query.from as string;
      const to = req.query.to as string;

      const filters = { title, author, from, to };
      const { artworks } = await this.artworkUseCase.list(filters);

      return res.json(artworks);
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Failed to retrieve artworks" });
    }
  }
}

