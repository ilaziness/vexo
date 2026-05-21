import React from "react";
import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Divider,
  Typography,
} from "@mui/material";

export type SettingsTab = "general" | "terminal" | "sync" | "ai" | "about";

interface NavItem {
  key: SettingsTab;
  label: string;
}

const navItems: NavItem[] = [
  { key: "general", label: "通用" },
  { key: "terminal", label: "终端" },
  { key: "sync", label: "同步" },
  { key: "ai", label: "AI" },
  { key: "about", label: "关于" },
];

interface SettingsNavProps {
  activeTab: SettingsTab;
  onTabChange: (tab: SettingsTab) => void;
}

const SettingsNav: React.FC<SettingsNavProps> = ({ activeTab, onTabChange }) => {
  return (
    <Box>
      <Box sx={{ p: 2, pb: 1.5 }}>
        <Typography variant="h6" sx={{ fontWeight: 600, fontSize: "1.2rem" }}>
          设置
        </Typography>
      </Box>
      <Divider />
      <List sx={{ py: 1 }}>
        {navItems.map((item) => (
          <ListItem key={item.key} disablePadding>
            <ListItemButton
              selected={activeTab === item.key}
              onClick={() => onTabChange(item.key)}
              sx={{
                py: 1,
                px: 2,
                "&.Mui-selected": {
                  backgroundColor: "primary.main",
                  color: "primary.contrastText",
                  "&:hover": {
                    backgroundColor: "primary.dark",
                  },
                },
              }}
            >
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );
};

export default SettingsNav;
