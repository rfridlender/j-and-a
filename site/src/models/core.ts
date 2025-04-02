import { z } from "zod"

export const coreSchema = z.object({
    createdAt: z.string(),
    createdBy: z.string().min(1),
    deletedAt: z.string(),
    deletedBy: z.string(),
})
