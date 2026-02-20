import { Request, Response } from "express";
import { UserUseCase } from "../useCases/UserUseCase";

export class UserController {
  constructor(private userUseCase: UserUseCase) {}

  async getProfile(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const id = parseInt(idParam);
      if (isNaN(id)) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const user = await this.userUseCase.getProfile(id);
      return res.json(user);
    } catch (error: any) {
      if (error.message === "User not found") {
        return res.status(404).json({ error: error.message });
      }
      return res.status(500).json({ error: "Error fetching user" });
    }
  }

  async updateProfile(req: Request, res: Response) {
    try {
      const userId = (req as any).user_id as number;
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const targetUserId = parseInt(idParam);
      if (isNaN(targetUserId)) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const { username, password } = req.body;

      const result = await this.userUseCase.updateProfile(userId, targetUserId, {
        username,
        password,
      });

      return res.json(result);
    } catch (error: any) {
      if (error.message === "You can only update your own profile") {
        return res.status(403).json({ error: error.message });
      }
      if (error.message === "User not found") {
        return res.status(404).json({ error: error.message });
      }
      return res.status(400).json({ error: "Invalid JSON" });
    }
  }

  async getUserArtworks(req: Request, res: Response) {
    try {
      const idParam = Array.isArray(req.params.id) ? req.params.id[0] : req.params.id;
      if (!idParam) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const userId = parseInt(idParam);
      if (isNaN(userId)) {
        return res.status(400).json({ error: "Invalid user ID" });
      }
      const artworks = await this.userUseCase.getUserArtworks(userId);
      return res.json(artworks);
    } catch (error: any) {
      console.error(error);
      return res.status(500).json({ error: "Error fetching artworks" });
    }
  }
}

