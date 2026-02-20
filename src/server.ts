import express from "express";
import cors from "cors";
import morgan from "morgan";
import path from "path";
import dotenv from "dotenv";

import { registerArtworkRoutes } from "./routes/artwork";
import { registerAuthRoutes } from "./routes/auth";
import { registerUserRoutes } from "./routes/users";
import { prisma } from "./database/prisma";

dotenv.config();

async function bootstrap() {
  // Testar conexão com Prisma
  try {
    await prisma.$connect();
    console.log("Database connected successfully");
  } catch (error) {
    console.error("Failed to connect to database:", error);
    process.exit(1);
  }

  const app = express();

  app.use(
    express.json({
      limit: "50mb",
    }),
  );
  app.use(
    express.urlencoded({
      extended: true,
      limit: "50mb",
    }),
  );

  app.use(morgan("dev"));
  app.use(
    cors({
      origin: "*",
    }),
  );

  app.get("/", (_req, res) => {
    res.send("Welcome to My Studio!");
  });

  registerArtworkRoutes(app);
  registerAuthRoutes(app);
  registerUserRoutes(app);

  const uploadsDir = path.join(__dirname, "..", "uploads");
  app.use("/uploads", express.static(uploadsDir));

  const PORT = process.env.PORT || 3000;
  app.listen(PORT, () => {
    console.log(`My Studio API listening on port ${PORT}`);
  });
}

bootstrap().catch((err) => {
  console.error("Failed to start server", err);
  process.exit(1);
});

// Graceful shutdown
process.on("SIGINT", async () => {
  await prisma.$disconnect();
  process.exit(0);
});

process.on("SIGTERM", async () => {
  await prisma.$disconnect();
  process.exit(0);
});
