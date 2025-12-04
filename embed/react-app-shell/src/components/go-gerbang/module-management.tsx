/* eslint-disable @typescript-eslint/no-explicit-any */

import { useState } from "react";
import { EllipsisVertical, Files } from "lucide-react";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { BackendUrlBase } from "@/services/baseService";

type DetailModuleProps = {
  name: string
  listModules?: {
    moduleName: string
    moduleType: string
    moduleLocation: string
  }[]
}

function DetailModule({
  name,
  listModules
}: DetailModuleProps) {
  const [openDialog, setOpenDialog] = useState(false);

  const [fileName, setFileName] = useState<string>("");
  const [formData] = useState(() => new FormData());

  // function sendForm(form: FormData) {
  async function sendForm() {
    const location = formData.get("file-location");
    // const file = formData.get("file-upload");
    const files = formData.get("files");
    // const formSend = {
    //   file_location: location,
    //   file: file,
    // }
    console.log(location);
    console.log(files);
    if (location && files) {
      const response = await fetch(`${BackendUrlBase}/upload-file`, {
        method: "POST",
        body: formData,
      });
      const result = await response.json();
      console.log("Response:", result);
    }
  }

  const handleDrop: React.DragEventHandler<HTMLDivElement> = async (e) => {
    e.preventDefault();
    // const files = e?.dataTransfer?.files;
    // if (files && files.length > 0) {
    //   setFormFile(files[0]);
    // }
    const items = e.dataTransfer.items;

    for (const item of items) {
      const entry = item.webkitGetAsEntry();
      if (!entry) continue;

      if (entry.isFile) {
        await readFileEntry(entry, ""); // root folder
      }

      if (entry.isDirectory) {
        await readDirectoryRecursive(entry, entry.name); // keep folder name
      }
    }
  }

  function readFileEntry(entry: any, path: string): Promise<void> {
    return new Promise((resolve) => {
      entry.file((file: File) => {
        const fullPath = path ? `${path}/${file.name}` : file.name;
        formData.append("files", file, fullPath);

        setFileName("folder");
        resolve();
      });
    });
  }

  function readDirectoryRecursive(entry: any, path: string): Promise<void> {
    return new Promise((resolve) => {
      const reader = entry.createReader();

      reader.readEntries(async (entries: any[]) => {
        for (const e of entries) {
          if (e.isFile) {
            await readFileEntry(e, path);
          } else if (e.isDirectory) {
            await readDirectoryRecursive(e, `${path}/${e.name}`);
          }
        }
        resolve();
      });
    });
  }

  const clearFile = () => {
    formData.delete("files");
    setFileName("");
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{name} MODULE</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Module Name</TableHead>
              <TableHead>Module Type</TableHead>
              <TableHead className="w-10"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listModules?.map((module) => (
              <TableRow key={module.moduleLocation}>
                <TableCell>{module.moduleName}</TableCell>
                <TableCell>{module.moduleType}</TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant="ghost"
                        className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
                        size="icon"
                      >
                        <EllipsisVertical />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-32">
                      <DropdownMenuItem onClick={() => setOpenDialog(true)}>Edit</DropdownMenuItem>
                      <DropdownMenuItem>Make a copy</DropdownMenuItem>
                      <DropdownMenuItem>Favorite</DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive">Delete</DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <Dialog open={openDialog} onOpenChange={setOpenDialog}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Module</DialogTitle>
              <DialogDescription>Here to upload new module file</DialogDescription>
            </DialogHeader>
            <form action={sendForm}>
              <Input type="text" placeholder="Location" onChange={(e) => formData.set("file-location", e.target.value)} required />
              <div className="mt-2 flex justify-center rounded-lg border-dashed border-2 px-6 py-10" onDrop={handleDrop} onDragOver={(e) => {e.preventDefault()}}>
                <div className="text-center">
                  <Files className="mx-auto size-10" />
                  <div className="mt-4 flex justify-center text-sm/6">
                    <Label htmlFor="file-upload">
                      <span>{fileName ? fileName : "Click here to Upload"}</span>
                      <input id="file-upload" type="file" name="file-upload" className="sr-only" {...({ webkitdirectory: "true", directory: "" } as any)} onChange={(e) => {
                        const files = e.target.files;
                        if (!files) return;

                        for (const file of files) {
                          formData.append("files", file, file.webkitRelativePath);
                        }
                      }} />
                    </Label>
                  </div>
                  {!fileName && <p className="text-sm mt-2 italic">or drag and drop here</p>}
                </div>
              </div>
              <input name="query" />
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline" onClick={clearFile}>Cancel</Button>
                </DialogClose>
                <Button type="submit">Upload</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  )
}


export default function ModuleuManagement() {
  const frontendModule = [
    {
      moduleName: "dashboard",
      moduleType: "React",
      moduleLocation: "/var/www/html/dashboard/"
    },
    {
      moduleName: "risk-register",
      moduleType: "React",
      moduleLocation: "/var/www/html/risk-register/"
    },
    {
      moduleName: "risk-timeline",
      moduleType: "React",
      moduleLocation: "/var/www/html/risk-timeline/"
    },
  ]


  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <DetailModule name="BACKEND" />
      </div>
      
      <div>
        <DetailModule name="FRONTEND" listModules={frontendModule} />
      </div>
    </div>
  )
}