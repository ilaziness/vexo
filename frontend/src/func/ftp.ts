import { FileInfo } from "./types";

export const sortFileList = (files: FileInfo[]): FileInfo[] => {
  return [...files].sort((a, b) => {
    if (a.isDir && !b.isDir) return -1;
    if (!a.isDir && b.isDir) return 1;
    return a.name.localeCompare(b.name, undefined, {
      numeric: true,
      sensitivity: "base",
    });
  });
};
