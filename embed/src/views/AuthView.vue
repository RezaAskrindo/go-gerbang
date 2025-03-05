<template>
  <div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
    <div class="w-full max-w-sm">
      <div class="flex flex-col gap-6">
        <RouterView />
        <div class="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 hover:[&_a]:text-primary">
          Silahkan Login atau Daftar Jika Belum Punya Akun
        </div>
      </div>
    </div>
  </div>
</template>