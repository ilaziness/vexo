## 1. Setup Dependencies

- [x] 1.1 Install @hello-pangea/dnd: `npm install @hello-pangea/dnd`

## 2. Store Enhancement

- [x] 2.1 Add `reorderTabs` method to `useSSHTabsStore` in `frontend/src/stores/ssh.ts`
- [x] 2.2 Implement array reordering logic using immutable update pattern

## 3. Sortable Tab Component

- [x] 3.1 Create draggable tab wrapper using `Draggable` from @hello-pangea/dnd
- [x] 3.2 Implement drag handle and visual feedback (opacity, shadow)
- [x] 3.3 Add CSS styles for dragging state

## 4. SSHTabs Integration

- [x] 4.1 Wrap Tabs with `DragDropContext` in `SSHTabs.tsx`
- [x] 4.2 Add `Droppable` with horizontal direction for the tab list
- [x] 4.3 Wrap each Tab with `Draggable` component
- [x] 4.4 Implement `onDragEnd` handler to call reorderTabs

## 5. Visual Polish

- [x] 5.1 Add drag overlay for smooth visual feedback during drag
- [x] 5.2 Style the drop indicator/placeholder
- [x] 5.3 Ensure scroll buttons still work with draggable tabs

## 6. Testing & Verification

- [x] 6.1 Test dragging tabs to reorder
- [x] 6.2 Verify SSH session state persists after reorder
- [x] 6.3 Test cancel drag with Escape key
- [x] 6.4 Test with multiple tabs and scroll behavior
