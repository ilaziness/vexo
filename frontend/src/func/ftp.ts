import { FileInfo } from "./types";

export const sortFileList = (files: FileInfo[]): FileInfo[] => {
  return files.sort((a, b) => {
    // 目录优先于文件
    if (a.isDir && !b.isDir) return -1;
    if (!a.isDir && b.isDir) return 1;
    // 目录和文件内部按名称自然排序
    return a.name.localeCompare(b.name, undefined, {
      numeric: true,
      sensitivity: "base",
    });
  });
};
