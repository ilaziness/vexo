import React from "react";
import { Browser } from "@wailsio/runtime";
import {
  Button,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
} from "@mui/material";
import Markdown from "markdown-to-jsx";
import { NewVersion } from "../../bindings/github.com/ilaziness/vexo/services";

interface UpdateAvailableDialogProps {
  open: boolean;
  onClose: () => void;
  newVersion: NewVersion | null;
}

function openExternalURL(url: string | undefined | null) {
  if (!url) {
    return;
  }
  try {
    const parsed = new URL(url);
    if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
      return;
    }
    Browser.OpenURL(parsed.toString());
  } catch {
    // ignore invalid URLs
  }
}

const markdownOptions = {
  disableParsingRawHTML: true,
  overrides: {
    h1: {
      component: Typography,
      props: { variant: "h6", gutterBottom: true, sx: { fontWeight: 700, mt: 1 } },
    },
    h2: {
      component: Typography,
      props: {
        variant: "subtitle1",
        gutterBottom: true,
        sx: { fontWeight: 700, mt: 1 },
      },
    },
    h3: {
      component: Typography,
      props: {
        variant: "subtitle2",
        gutterBottom: true,
        sx: { fontWeight: 600, mt: 1 },
      },
    },
    p: {
      component: Typography,
      props: { variant: "body2", paragraph: true },
    },
    li: {
      component: Typography,
      props: { component: "li", variant: "body2" },
    },
    a: {
      component: ({
        href,
        children,
        ...rest
      }: React.AnchorHTMLAttributes<HTMLAnchorElement>) => (
        <Box
          component="a"
          href={href}
          {...rest}
          onClick={(e: React.MouseEvent) => {
            e.preventDefault();
            openExternalURL(href);
          }}
          sx={{
            color: "primary.main",
            textDecoration: "underline",
            cursor: "pointer",
          }}
        >
          {children}
        </Box>
      ),
    },
    code: {
      component: ({ children }: { children?: React.ReactNode }) => (
        <Box
          component="code"
          sx={{
            fontFamily: "monospace",
            fontSize: "0.85em",
            px: 0.5,
            py: 0.25,
            borderRadius: 0.5,
            bgcolor: "action.hover",
          }}
        >
          {children}
        </Box>
      ),
    },
    pre: {
      component: ({ children }: { children?: React.ReactNode }) => (
        <Box
          component="pre"
          sx={{
            fontFamily: "monospace",
            fontSize: "0.8rem",
            p: 1,
            borderRadius: 1,
            bgcolor: "action.hover",
            overflow: "auto",
            my: 1,
          }}
        >
          {children}
        </Box>
      ),
    },
  },
};

export default function UpdateAvailableDialog({
  open,
  onClose,
  newVersion,
}: UpdateAvailableDialogProps) {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>发现新版本</DialogTitle>
      <DialogContent>
        <Typography sx={{ fontWeight: 600 }}>{newVersion?.Version}</Typography>
        {newVersion?.Notes ? (
          <Box
            sx={{
              mt: 1,
              maxHeight: 360,
              overflow: "auto",
              "& ul, & ol": { pl: 2.5, m: 0 },
            }}
          >
            <Markdown options={markdownOptions}>{newVersion.Notes}</Markdown>
          </Box>
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>关闭</Button>
        <Button onClick={() => openExternalURL(newVersion?.URL)}>
          打开下载页
        </Button>
      </DialogActions>
    </Dialog>
  );
}
