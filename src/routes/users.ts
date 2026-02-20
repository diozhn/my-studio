import { Router } from "express";
import { UserController } from "../controllers/UserController";
import { UserUseCase } from "../useCases/UserUseCase";
import { UserRepository } from "../repositories/UserRepository";
import { ArtworkRepository } from "../repositories/ArtworkRepository";
import { authMiddleware } from "./auth";

const router = Router();

// Inicializar dependências
const userRepository = new UserRepository();
const artworkRepository = new ArtworkRepository();
const userUseCase = new UserUseCase(userRepository, artworkRepository);
const userController = new UserController(userUseCase);

export function registerUserRoutes(app: import("express").Express) {
  app.use(router);

  router.get("/users/:id", authMiddleware, (req, res) =>
    userController.getProfile(req, res),
  );
  router.patch("/users/:id", authMiddleware, (req, res) =>
    userController.updateProfile(req, res),
  );
  router.get("/users/:id/artworks", (req, res) =>
    userController.getUserArtworks(req, res),
  );
}
