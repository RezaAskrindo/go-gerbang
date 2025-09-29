
export interface AppTobarProps {
  logoName?: string
  logoImage?: string
  labels?: string[]
}

export interface AppLayoutProps extends AppTobarProps {
  layoutType?: string;
}