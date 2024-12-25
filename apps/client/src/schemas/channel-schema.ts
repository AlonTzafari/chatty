import { z } from "zod";

const channelSchema = z.object({
    Id: z.string(), 
    Name: z.string(),
    CreatedAt: z.string(), 
    Avatar: z.string(),
})

export default channelSchema