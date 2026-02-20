import { Router } from "express";
import { AuthController } from "../controllers/AuthController";
import { AuthUseCase } from "../useCases/AuthUseCase";
import { UserRepository } from "../repositories/UserRepository";
import { requireAuth } from "../middlewares/auth";

const router = Router();

// Inicializar dependências
const userRepository = new UserRepository();
const authUseCase = new AuthUseCase(userRepository);
const authController = new AuthController(authUseCase);

// Exportar requireAuth para uso em outras rotas
export const authMiddleware = requireAuth(authUseCase);

export function registerAuthRoutes(app: import("express").Express) {
  app.use(router);

  router.post("/register", (req, res) => authController.register(req, res));
  router.post("/login", (req, res) => authController.login(req, res));
  router.post("/refresh-token", (req, res) => authController.refreshToken(req, res));
}
