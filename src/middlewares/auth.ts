import { Request, Response, NextFunction } from "express";
import { AuthUseCase } from "../useCases/AuthUseCase";

export function requireAuth(authUseCase: AuthUseCase) {
  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      const authHeader = req.headers.authorization;

      if (!authHeader || !authHeader.startsWith("Bearer ")) {
        return res.status(401).json({ error: "Token not provided" });
      }

      const token = authHeader.substring(7);
      const decoded = authUseCase.verifyToken(token);

      (req as any).user_id = decoded.user_id;
      next();
    } catch (error: any) {
      return res.status(401).json({ error: "Invalid Token" });
    }
  };
}

