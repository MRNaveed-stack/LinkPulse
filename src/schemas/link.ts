import { z } from 'zod';

export const createLinkSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  slug: z
    .string()
    .min(1, 'Slug is required')
    .regex(/^[a-zA-Z0-9-]+$/, 'Only letters, numbers, and hyphens allowed'),
  destination_url: z.url('Must be a valid URL'),
});

// For update, all fields are optional but must pass same rules if provided
export const updateLinkSchema = z.object({
  title: z.string().min(1).optional(),
  slug: z
    .string()
    .min(1)
    .regex(/^[a-zA-Z0-9-]+$/)
    .optional(),
  destination_url: z.url().optional(),
});

export type CreateLinkInput = z.infer<typeof createLinkSchema>;
export type UpdateLinkInput = z.infer<typeof updateLinkSchema>;
