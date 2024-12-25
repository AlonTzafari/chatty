import { z } from "zod";

const authSchema = z.object({
    Id: z.string(),
    Username: z.string()
})

export default authSchema