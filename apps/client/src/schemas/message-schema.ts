import { z } from "zod";

const messageSchema = z.object({
    Id: z.string(), 
    UserId: z.string(),
    Content: z.string(), 
    CreatedAt: z.string(), 
    Username: z.string(),
})

export default messageSchema