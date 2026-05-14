import { useState } from "react";
import { Plus } from "lucide-react";

import { Button } from "@/components/ui/button";

import SheetForm from "@/components/sheet-form";
import CardInformation from "@/components/card-information";
import useSWR from "swr";
import { fetchSWR } from "@/services/use-swr-service";
import { BackendUrlBase } from "@/services/baseService";

export default function CaddyManagement() {
  const [openSheet, setOpenSheet] = useState(false);

  const { data, isLoading } = useSWR(`${BackendUrlBase}/check-local-service?url=http://localhost:2019/config&getRes=true`, fetchSWR);

  const [dataForm, setDataFrom] = useState();

  console.log(data)

  return (
    <div className="flex-1 flex-col gap-8 md:flex">
      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-col gap-1">
          <h2 className="text-2xl font-semibold tracking-tight">
            Caddy Management
          </h2>
          <p className="text-muted-foreground">
            Here&apos;s list Caddy Setup
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={() => setOpenSheet(true)} variant="outline">
            <Plus />
            Setup
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 gap-4">
        <div>
          <CardInformation name="Resend Setup" description="There is list of Resend Email for Notification" />
        </div>
      </div>
      <SheetForm 
        name="Caddy Form"
        openSheet={openSheet} 
        setOpenSheet={setOpenSheet}
      >
        <div>Caddy Form</div>
      </SheetForm>
      {/* <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete {dataForm?.sender} Notification?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will delete the data. Data deleted cannot be restored.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogAction>Cancel</AlertDialogAction>
            <AlertDialogCancel variant="destructive" onClick={() => {
              if (dataForm?.notif_type) {
                toast.promise(
                  useDeleteConfiguration(dataForm?.notif_type, dataForm?.sender),
                  {
                    loading: "Waiting...",
                    success: () => {
                      mutateData();
                      setOpenDialog(false);
                      return "Success"
                    },
                    error: (err) => {
                      setOpenDialog(false);
                      return err.message || "Failed"
                    },
                  } 
                )
              }
            }}>Delete</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog> */}
    </div>
  )
}