# RFC: smpatch V3 设计文档

## 一、设计目标

为任意 `map[string]any` 结构提供一套**显式声明的通用 Strategic Merge Patch** 能力，不依赖 OpenAPI schema 或 Go struct tag，通过 PathKey + Ope 显式声明操作类型和匹配方式。

---

## 二、核心概念

### 2.1 PathKey

| 规则 | 说明 |
|---|---|
| 不以 `/` 开头 | `spec/bindings/role=="admin"/members` |
| 为空字符串 `""` | ❌ 不支持，视为非法（根路径操作被拒绝） |
| `/` 分隔符 | 每一段称为一个 segment |
| 普通 segment | map 字段名，如 `spec`、`bindings` |
| 表达式 segment | 含 `==`、`!=`、`&&`、`||` 的 segment，用于筛选 struct 数组元素 |

### 2.2 Patch 结构体

```go
type Patch struct {
    Ope      string   // merge / replace / delete
    PathKey  string   // 路径，如 "spec/bindings"
    Value    any      // 操作值
    ByKey    string   // struct 数组的合并键（mapAr 使用）
    ItemOps  string   // 标量数组操作类型（itemOps 使用）
    MixedAr  bool     // 是否走 mixed 全量替换
    Old      any      // itemOps replace 时的旧值
}
```

---

## 三、操作类型总览

| Ope | 适用目标 | 说明 |
|---|---|---|
| `merge` | map / struct 数组 / 标量数组 | 根据附加字段决定具体策略 |
| `replace` | map / 任意数组 | 全量替换目标字段 |
| `delete` | map / 数组元素 | 删除指定路径或匹配的元素 |

---

## 四、数组处理策略（核心设计）

YAML/JSON 中数组没有标准的 merge 语义，因此按数组元素类型拆分为三种策略：

### 4.1 struct 数组 → `mapAr` + `ByKey`

- **定位方式**：PathKey 普通路径走到数组字段
- **匹配方式**：`ByKey` 等值匹配（如 `role` 字段）
- **操作语义**：对整个匹配的 struct 做增删改
- **与 K8s 的关系**：等价于 `listType: map` + `listKey: role`

```
PathKey: "spec/bindings"
ByKey:   "role"
Value:   [{role: "admin", members: ["3"]}]
```

→ 找到 `role == "admin"` 的 binding，将 Value 字段 merge 进去

### 4.2 标量数组 → `itemOps`

- **定位方式**：PathKey 含表达式 segment，筛选父级 struct
- **匹配方式**：表达式求值（如 `role=="admin"`），必须唯一匹配 1 个元素
- **操作目标**：匹配 struct 内部的某个数组字段
- **操作类型**：`add` / `remove` / `replace` / `keep` / `disable`
- **去重**：`add` 时自动去重（`slices.Contains`）
- **与 K8s 的关系**：等价于 `listType: set`

```
PathKey:  "spec/bindings/role==\"admin\"/members"
ItemOps:  "add"
Value:    ["3"]
```

→ 找到 `role == "admin"` 的 binding，向其 `members` 数组追加 `"3"`

### 4.3 mixed 数组 → `mixed`

- **定位方式**：PathKey 普通路径走到数组字段
- **操作语义**：全量替换，不做任何匹配或去重
- **去重责任**：由用户保证
- **与 K8s 的关系**：等价于 `listType: atomic`

```
PathKey: "spec/p"
MixedAr: true
Value:   [321, 987]
```

→ `p` 数组被整体替换为 `[321, 987]`

---

## 五、普通 map 合并 → `dpMeg`

- **定位方式**：PathKey 普通路径走到目标 map
- **操作语义**：深度拷贝 + 逐 key 覆盖（已有的 key 被覆盖，没有的保留）
- **nil 保护**：中间路径不存在时自动初始化为空 map

```
PathKey: "spec"
Value:   {add: "r"}
```

→ `spec.add` 被设置为 `"r"`，`spec.bindings` 等其他字段不受影响

---

## 六、入口路由 → `Apply`

所有操作统一入口为 `Apply(src, patches, tgt)`：

```
┌─────────────────────────────────────────────────────┐
│              RFC V3 设计矩阵                         │
├──────────────┬───────────┬──────────────────────────┤
│ 数组类型      │ 操作       │ 代码路径                 │
├──────────────┼───────────┼──────────────────────────┤
│ struct       │ merge     │ apply.go → mapAr         │ ✅
│ struct       │ delete    │ apply.go → applyDel      │ ✅
│ scalar       │ add/rem/  │ apply.go → itemOps       │ ✅
│              │ rep/keep/ │ (表达式筛选)              │
│              │ disable   │                          │
│ mixed        │ replace   │ apply.go → mixed         │ ✅
│ map(非数组)  │ merge     │ apply.go → dpMeg         │ ✅
│ map(非数组)  │ replace   │ apply.go → applyRpl      │ ✅
│ 任意         │ delete    │ apply.go → applyDel      │ ✅
└──────────────┴───────────┴──────────────────────────┘
```

### switch 匹配顺序（从高到低）

1. `delete` — 无条件优先
2. `replace` + `MixedAr:true` — mixed 全量替换
3. `replace` — 普通全量替换
4. `merge` + `ByKey != ""` — struct 数组合并
5. `merge` + `ItemOps != ""` — 标量数组操作
6. `merge` — 普通 map 深度合并
7. `default` — 报错

### 入口熔断

- PathKey 为空（`""` 或 `"/"` trim 后为空）→ 直接返回 error，拒绝根路径操作

---

## 七、与 Kustomize Strategic Merge Patch 语义对照

### 设计动机

Kustomize 的 Strategic Merge Patch（SMP）对内置资源（如 `Deployment.spec.containers`）通过 Go struct tag 或 OpenAPI schema 隐式确定数组合并策略。但对 Custom Resource（CR）默认**全量替换**，需要额外声明 schema 才能按 key 合并。

`smpatch` 的设计目标是：**不依赖 schema，通过显式声明操作类型和匹配方式，为任意 CR 提供与 Kustomize SMP 等价甚至更灵活的数组处理能力。**

### 语义映射表

| Kustomize SMP 行为 | 触发条件 | smpatch 对应策略 | 匹配/定位方式 |
|---|---|---|---|
| struct 数组按 mergeKey 合并 | CRD 声明 `patchStrategy:"merge"` + `patchMergeKey:"role"` | `mapAr` + `ByKey` | ByKey 等值匹配 |
| struct 数组全量替换 | 无 schema 或未声明 merge | `mixed`（全量替换） | 无匹配，直接覆盖 |
| 标量数组集合操作 | `patchStrategy:"merge"` 声明为 set | `itemOps`（add/remove/replace/keep） | 表达式筛选父级 struct → 操作内部数组字段 |
| 普通 map 字段合并 | 任意 | `dpMeg`（深度合并） | PathKey 普通路径 |
| JSON Patch 删除 | `patchesJson6902` | `applyDel` | PathKey 定位 |

### 核心差异

| 维度 | Kustomize SMP | smpatch |
|---|---|---|
| schema 依赖 | 需要 OpenAPI schema 或 Go struct tag | ❌ 不需要，PathKey + Ope 显式声明 |
| 表达式筛选 | ❌ 不支持（只能按 mergeKey 等值匹配） | ✅ 支持（`role=="admin" && env=="prod"`） |
| 标量数组操作 | 有限（需声明为 set） | ✅ 完整 CRUD（add/remove/replace/keep） |
| 适用范围 | K8s 资源 | 任意 `map[string]any` 结构 |
| 接入方式 | kustomize binary 或 krusty Go API | 纯 Go 库，直接调用 `Apply()` |

### 总结

> **smpatch 本质上是一套"显式声明的通用 Strategic Merge Patch"：**
> - `mapAr` + `ByKey` ↔ Kustomize 的 `patchMergeKey` 合并
> - `mixed` ↔ Kustomize 默认的全量替换
> - `itemOps` ↔ Kustomize 缺失的标量数组精细操作
> - `dpMeg` ↔ Kustomize 的 map 字段深度合并

---

## 八、代码文件职责划分

| 文件 | 职责 |
|---|---|
| `apply.go` | 统一入口 `Apply()`，路由分发，PathKey 空值熔断 |
| `dpmeg.go` | 普通 map 深度合并，nil 路径自动初始化 |
| `mapar.go` | struct 数组按 ByKey 匹配，整体增删改 |
| `itemops.go` | walk 路径遍历（含表达式 segment 识别） |
| `exprl.go` | expr 表达式编译、执行、唯一性校验 |
| `ictrler.go` | itemOps 最终操作执行（add/remove/replace/keep） |
| `mixed.go` | mixed 数组全量替换 |
| `applydel.go` | 删除操作 |
| `applyrpl.go` | 普通 replace 操作 |
| `types.go` | Patch 结构体定义 |
| `diffcmp.go` | 差异比较工具 |
| `mapar.go` 中的通用辅助 | `DeepCopy`、`Contains` 等 |

---

## 九、设计约束与边界

| 约束 | 说明 |
|---|---|
| PathKey 不能为空 | 根路径操作不支持，入口拒绝 |
| 表达式必须唯一匹配 | itemOps 中表达式匹配 0 个或多个均报错 |
| 表达式 segment 前必须是数组 | 否则报错 |
| 最终目标必须是数组 | itemOps 操作的字段必须是 `[]any` |
| mapAr 的 ByKey 不能为空 | 否则报错 |
| 标量数组 add 自动去重 | `slices.Contains` 判断 |
| mixed 不去重 | 压力交给用户 |
| 不支持跨层级引用 | 如 `bindings[0].role` |
| 不支持函数调用 | 如 `contains(role, "adm")` |
| 不支持正则匹配 | 如 `role=~"adm.*"` |