import { type ReactNode } from "react";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

type TSheetForm = {
  name: string
  children: ReactNode
  openSheet: boolean
  setOpenSheet: (v: boolean) => void
  description?: string
}

export default function SheetForm({
  name,
  children,
  description,
  openSheet,
  setOpenSheet,
}: TSheetForm) {

  return (
    <Sheet open={openSheet} onOpenChange={setOpenSheet}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>{name}</SheetTitle>
          <SheetDescription className={description || "sr-only"}>{description ? description : "no description"}</SheetDescription>
        </SheetHeader>
        { children }
      </SheetContent>
    </Sheet>
  )
}
