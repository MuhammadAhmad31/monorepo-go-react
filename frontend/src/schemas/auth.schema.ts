import { z } from "zod"
import type { LoginRequest } from "@/generated/models"

export const loginSchema = z.object({
  email: z.email("Email tidak valid"),
  password: z.string().min(6, "Password minimal 6 karakter"),
}) satisfies z.ZodType<LoginRequest>
