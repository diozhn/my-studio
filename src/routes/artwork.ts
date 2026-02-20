import { Router } from "express";
import { ArtworkController, upload } from "../controllers/ArtworkController";
import { ArtworkUseCase } from "../useCases/ArtworkUseCase";
import { ArtworkRepository } from "../repositories/ArtworkRepository";
import { authMiddleware } from "./auth";

const router = Router();

// Inicializar dependências
const artworkRepository = new ArtworkRepository();
const artworkUseCase = new ArtworkUseCase(artworkRepository);
const artworkController = new ArtworkController(artworkUseCase);

export function registerArtworkRoutes(app: import("express").Express) {
  app.use(router);

  router.post("/artworks", authMiddleware, upload.single("image"), (req, res) =>
    artworkController.create(req, res),
  );
  router.post("/artworks/:id/like", (req, res) => artworkController.like(req, res));
  router.get("/gallery", (req, res) => artworkController.getGallery(req, res));
  router.get("/top-artworks", (req, res) => artworkController.getTopArtworks(req, res));
  router.get("/artworks", (req, res) => artworkController.list(req, res));
  router.get("/artworks/:id", (req, res) => artworkController.getById(req, res));
  router.get("/artworks/filter", (req, res) => artworkController.getFiltered(req, res));
  router.patch("/artworks/:id", authMiddleware, (req, res) =>
    artworkController.update(req, res),
  );
  router.delete("/artworks/:id", authMiddleware, (req, res) =>
    artworkController.delete(req, res),
  );
}
