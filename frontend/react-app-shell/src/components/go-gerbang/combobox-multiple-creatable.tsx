"use client"

import * as React from "react"
import {
  Combobox,
  ComboboxChip,
  ComboboxChips,
  ComboboxChipsInput,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxItem,
  ComboboxList,
  ComboboxValue,
  useComboboxAnchor,
} from "@/components/ui/combobox"
import { cn } from "@/lib/utils"

const initialData: string[] = [] as const

interface LabelItem {
  creatable?: string
  id: string
  value: string
}

interface ComboboxMultipleCreatableProps {
  className?: string
  placeholder?: string
  value?: string[]
  onValueChange?: (value: string[]) => void
}

export function ComboboxMultipleCreatable({
  className,
  placeholder,
  value: controlledValue,
  onValueChange: controlledOnValueChange,
}: ComboboxMultipleCreatableProps) {
  const anchor = useComboboxAnchor()
  const id = React.useId()

  const [labels, setLabels] = React.useState<LabelItem[]>(
    initialData.map((fw) => ({
      id: fw.toLowerCase().replace(/\s+/g, "-"),
      value: fw,
    }))
  )

  // Uncontrolled fallback state as strings if not used with React Hook Form
  const [internalValue, setInternalValue] = React.useState<string[]>([])

  const selectedStrings = controlledValue ?? internalValue
  const setSelectedStrings = controlledOnValueChange ?? setInternalValue

  // Convert selected strings into LabelItem objects for Base UI Combobox compatibility
  const selectedItems: LabelItem[] = React.useMemo(() => {
    return selectedStrings.map((str) => {
      const existing = labels.find(
        (l) => l.value.trim().toLocaleLowerCase() === str.trim().toLocaleLowerCase()
      )
      return (
        existing ?? {
          id: str.toLowerCase().replace(/\s+/g, "-"),
          value: str,
        }
      )
    })
  }, [selectedStrings, labels])

  const [query, setQuery] = React.useState("")
  const highlightedItemRef = React.useRef<LabelItem | undefined>(undefined)

  function handleCreateNew(rawValue: string) {
    const trimmedValue = rawValue.trim()
    if (!trimmedValue) return

    const normalized = trimmedValue.toLocaleLowerCase()
    const baseId = normalized.replace(/\s+/g, "-")

    const existing = labels.find(
      (l) => l.value.trim().toLocaleLowerCase() === normalized
    )

    if (existing) {
      if (!selectedStrings.some((s) => s.trim().toLocaleLowerCase() === normalized)) {
        setSelectedStrings([...selectedStrings, existing.value])
      }
      setQuery("")
      return
    }

    const existingIds = new Set(labels.map((l) => l.id))
    let uniqueId = baseId
    if (existingIds.has(uniqueId)) {
      let i = 2
      while (existingIds.has(`${baseId}-${i}`)) {
        i += 1
      }
      uniqueId = `${baseId}-${i}`
    }

    const newItem: LabelItem = { id: uniqueId, value: trimmedValue }

    setLabels((prev) => [...prev, newItem])
    setSelectedStrings([...selectedStrings, newItem.value])
    setQuery("")
  }

  function handleInputKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key !== "Enter" || highlightedItemRef.current) {
      return
    }

    if (query.trim() !== "") {
      event.preventDefault()
      handleCreateNew(query)
    }
  }

  const trimmed = query.trim()
  const lowered = trimmed.toLocaleLowerCase()
  const exactExists = labels.some(
    (l) => l.value.trim().toLocaleLowerCase() === lowered
  )

  const itemsForView: Array<LabelItem> =
    trimmed !== "" && !exactExists
      ? [
          ...labels,
          {
            creatable: trimmed,
            id: `create:${lowered}`,
            value: `Create "${trimmed}"`,
          },
        ]
      : labels

  return (
    <Combobox
      multiple
      autoHighlight
      items={itemsForView}
      value={selectedItems}
      onValueChange={(next) => {
        const creatableSelection = next.find(
          (item) =>
            item.creatable &&
            !selectedItems.some((current) => current.value.trim().toLocaleLowerCase() === item.creatable?.trim().toLocaleLowerCase())
        )

        if (creatableSelection && creatableSelection.creatable) {
          handleCreateNew(creatableSelection.creatable)
          return
        }

        const cleanItems = next.filter((i) => !i.creatable)
        setSelectedStrings(cleanItems.map((i) => i.value))
        setQuery("")
      }}
      inputValue={query}
      onInputValueChange={setQuery}
      onItemHighlighted={(item) => {
        highlightedItemRef.current = item
      }}
    >
      <ComboboxChips ref={anchor} className={cn("w-full", className)}>
        <ComboboxValue>
          {(values: LabelItem[]) => (
            <React.Fragment>
              {values.map((label) => (
                <ComboboxChip key={label.id}>{label.value}</ComboboxChip>
              ))}
              <ComboboxChipsInput
                id={id}
                placeholder={values.length > 0 ? "" : (placeholder ?? "")}
                onKeyDown={handleInputKeyDown}
              />
            </React.Fragment>
          )}
        </ComboboxValue>
      </ComboboxChips>
      <ComboboxContent anchor={anchor}>
        <ComboboxEmpty>No {placeholder ?? ""} found.</ComboboxEmpty>
        <ComboboxList>
          {(item: LabelItem) =>
            item.creatable ? (
              <ComboboxItem key={item.id} value={item} className="text-primary font-medium">
                + Create &quot;{item.creatable}&quot;
              </ComboboxItem>
            ) : (
              <ComboboxItem key={item.id} value={item}>
                {item.value}
              </ComboboxItem>
            )
          }
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  )
}