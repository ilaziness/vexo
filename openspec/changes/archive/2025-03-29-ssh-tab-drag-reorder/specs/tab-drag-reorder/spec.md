## ADDED Requirements

### Requirement: Tab drag initiation
用户 SHALL 能够通过鼠标拖拽 SSH 标签页来启动排序操作。

#### Scenario: Start dragging a tab
- **WHEN** 用户在标签页上按下鼠标左键并开始移动
- **THEN** 系统进入拖拽状态并显示拖拽视觉反馈

### Requirement: Visual feedback during drag
拖拽过程中 SHALL 提供清晰的视觉反馈，包括被拖拽标签的半透明效果和占位符指示。

#### Scenario: Visual feedback while dragging
- **WHEN** 用户拖拽标签页时
- **THEN** 被拖拽的标签页显示为半透明状态
- **AND** 目标位置显示占位符指示

### Requirement: Tab reorder on drop
当用户释放鼠标完成拖拽时，SHALL 根据放置位置重新排列标签页顺序。

#### Scenario: Drop tab to new position
- **WHEN** 用户将标签页拖拽到新位置并释放鼠标
- **THEN** 标签页列表按照新的顺序重新排列
- **AND** 当前激活的标签保持不变

### Requirement: Cancel drag on escape
用户 SHALL 能够通过按 Escape 键取消正在进行的拖拽操作。

#### Scenario: Cancel drag with Escape key
- **WHEN** 用户正在拖拽标签页时按下 Escape 键
- **THEN** 拖拽操作被取消
- **AND** 标签页顺序恢复到拖拽前的状态

### Requirement: Maintain tab state after reorder
标签页重排序后 SHALL 保持原有的会话状态和内容不变。

#### Scenario: State persistence after reorder
- **WHEN** 标签页被重新排序后
- **THEN** 每个标签页对应的 SSH 会话保持连接状态
- **AND** 终端内容保持不变
