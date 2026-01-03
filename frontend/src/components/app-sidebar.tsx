"use client";

import * as React from "react";

import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import { TeamSwitcher } from "@/components/team-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useAuth } from "@/hooks/use-auth";
import type { Role } from "@/generated/rbac";
import {
  filterNavigationByPermissions,
  NAVIGATION,
} from "@/constants/navigation";
import { GalleryVerticalEnd } from "lucide-react";
import { ENV } from "@/constants/env";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { getCurrentUser } = useAuth();
  const { user, isLoading } = getCurrentUser();

  const userRole: Role = (user?.role as Role) || "user";

  const filteredNavigation = React.useMemo(() => {
    if (!user) return [];
    return filterNavigationByPermissions(NAVIGATION, userRole);
  }, [user, userRole]);

  const userData = {
    name: user?.email || "Guest",
    email: user?.email || "guest@example.com",
    avatar: "/avatars/default.jpg",
  };

  const team = {
    name: "Acme Inc",
    logo: GalleryVerticalEnd,
  };

  if (isLoading) {
    return (
      <Sidebar collapsible="icon" {...props}>
        <SidebarHeader>
          <div className="h-10 bg-gray-200 animate-pulse rounded" />
        </SidebarHeader>
        <SidebarContent>
          <div className="space-y-2 p-2">
            <div className="h-8 bg-gray-200 animate-pulse rounded" />
            <div className="h-8 bg-gray-200 animate-pulse rounded" />
            <div className="h-8 bg-gray-200 animate-pulse rounded" />
          </div>
        </SidebarContent>
      </Sidebar>
    );
  }

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher team={team} />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={filteredNavigation} />
        {ENV.ENVIRONMENT === "development" && (
          <div className="p-2 text-xs text-muted-foreground border-t mt-4">
            <p>
              Role: <strong>{userRole}</strong>
            </p>
            <p>
              Visible items: <strong>{filteredNavigation.length}</strong>
            </p>
            <p className="text-green-600">✓ Filtered by RBAC permissions</p>
          </div>
        )}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={userData} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
