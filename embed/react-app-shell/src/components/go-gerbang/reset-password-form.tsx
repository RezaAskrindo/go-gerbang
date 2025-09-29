import { z } from "zod";

import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "../ui/input"
import { useState } from "react"
import { BackendUrlBase, FetchCsrfToken } from "@/services/baseService"
import { toast } from "sonner";

const emailSchema = z.email({ pattern: z.regexes.html5Email });

async function fetchEmail(email: string): Promise<Response> {
  const csrfToken = await FetchCsrfToken();

  const url = `${BackendUrlBase}/api/v1/auth/request-reset-password?baseUrl=${BackendUrlBase}`;
  return fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      'X-SGCsrf-Token': csrfToken
    },
    body: JSON.stringify({
      identity: email,
    }),
  });
}

const ResetPasswordForm = () => {
  const [open, setOpen] = useState(false);
  const [loadingButton, setLoadingButton] = useState(false);
  const [email, setEmail] = useState("");

  const sendEmail = async () => {
    const validation = emailSchema.safeParse(email);
    if (!validation.success) {
      toast.error(validation.error.issues[0].message)
      // You could set an error state here for the UI
      return;
    }

    setLoadingButton(true);
    let response = await fetchEmail(email);

    if (!response.ok) {
      try {
        const errorData = await response.json();
        if (errorData?.message === "CSRF validation failed") {
          response = await fetchEmail(email);
        } else {
          setLoadingButton(false);
        }
      } catch {
        toast.error("Something Error, please Call Administrator");
      }
    } else {
      const result = await response.json();
      toast.success(result.message);
      window.location.reload();
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <a href="#" className="ml-auto text-sm underline-offset-4 hover:underline">
          Forgot your password?
        </a>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Forgot Your Password</DialogTitle>
          <DialogDescription>
            Fill your email to get link Reset Password
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4">
          <Input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="name@example.com" />
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline" className="me-auto">Cancel</Button>
          </DialogClose>
          { loadingButton ? <Button disabled={true}>
            Loading...
          </Button> : <Button onClick={sendEmail}>Send Email</Button> }
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default ResetPasswordForm