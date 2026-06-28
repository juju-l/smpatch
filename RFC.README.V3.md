# RFC：通用资源 Patch CRD（v3）

> **文档版本**：v3.0  
> **状态**：Draft  
> **适用范围**：自研 CRD / Kubernetes 资源  
> **设计目标**：语义精确、安全可控、GitOps 友好

---

## 1. 设计目标

在 v2 基础上新增：

- ✅ 支持 **scalar 数组禁用特定值**
- ✅ 明确 **disable 与 remove 的差异**
- ✅ 强化 **安全 / 合规场景表达能力**

---

## 2. Patch CRD 示例（完整）

```yaml
apiVersion: patch.mycompany.io/v1alpha1
kind: ResourcePatch
metadata:
  name: security-patch
spec:
  namespace: default
  targetRef:
    apiVersion: mycompany.io/v1alpha1
    kind: MyApp
    name: demo
  patches:
    # ✅ struct 数组修改
    - ope: merge
      pathKey: /spec/components
      byKey: name
      value:
        - name: web
          image: nginx:1.25

    # ✅ scalar 数组替换
    - ope: merge
      pathKey: /spec/tags
      itemOps: replace
      old: legacy
      value: production

    # ✅ scalar 数组禁用（新增）
    - ope: merge
      pathKey: /spec/admins
      itemOps: disable
      value:
        - alice@evil.com
        - root@internal
```

---

## 3. Patch 条目字段定义（v3）

```yaml
patches:
  - ope: merge | replace | delete
    pathKey: string
    byKey: string
    itemOps: add | remove | keep | replace | disable
    mixedAr: bool
    value: any //必须是数组，就算只有一个，也需要是数组的格式
    old: any
```

---

## 4. itemOps 行为全景

| itemOps | 语义 | 是否允许 value 为数组 | 典型场景 |
|-------|----|----|----|
| add | 不存在则加入 | ✅ | 标签 / 成员 |
| remove | 存在则删除 | ✅ | 清理 |
| replace | 精确替换 | ✅ | 值变更 |
| keep | 仅保留 | ✅ | 安全基线 |
| **disable** | **禁止存在** | ✅ | **合规 / 封禁** |

---

## 5. disable 语义定义（核心）

### ✅ 行为规则

```
for each item in value:
    if item exists in target array:
        remove it
```

- ✅ 幂等
- ✅ 不抛错
- ✅ 不依赖顺序

---

### ✅ disable vs remove（重要区别）

| 维度 | remove | disable |
|----|----|----|
| 语义 | 删除指定值 | **禁止该值出现** |
| GitOps diff | 中性 | **强安全信号** |
| 多次 apply | 稳定 | 稳定 |
| 合规表达 | ❌ | ✅ |

✅ **disable 更适合策略类 patch**

---

## 6. 完整 scalar 数组示例

### 原始

```yaml
admins:
  - alice@company.com
  - bob@company.com
  - root@internal
```

### Patch

```yaml
- ope: merge
  pathKey: /spec/admins
  itemOps: disable
  value:
    - root@internal
```

### 结果

```yaml
admins:
  - alice@company.com
  - bob@company.com
```

---

## 7. 行为决策矩阵（v3）

| 数组类型 | add | remove | replace | keep | disable | 全量替换 |
|--------|----|----|----|----|----|----|
| struct | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ |
| scalar | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| mixed | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## 8. Controller 校验规则（新增）

```
if itemOps == "disable":
    require(value != null)
    forbid(old)

if itemOps == "replace":
    require(old)

if mixedAr == true:
    require(ope == "replace")
```

---

## 9. 最佳实践（安全视角）

- ✅ 使用 `disable` 表达 **安全策略**
- ✅ 使用 `keep` 表达 **最小权限**
- ✅ 使用 `replace` 表达 **业务变更**
- ✅ 禁止在 disable patch 中使用 comment 以外说明

```yaml
itemOps: disable
comment: "PCI-DSS compliance requirement"
```

---

## 10. Status 设计

```yaml
status:
  type: object
  observedGeneration: int64
  conditions:
    - type: Applied
      status: "True"
      reason: ""
      message: ""
```

---

## 11. 下一步建议

当前方案已具备：

- ✅ 语义完备
- ✅ 安全可控
- ✅ GitOps 友好
- ✅ Controller 可实现

下一步可选：

- ✅ OpenAPI / CRD YAML（含 validation）
- ✅ controller-runtime Controller（Go）
- ✅ Policy / Security Patch 示例集
- ✅ 架构评审 PPT

如需继续推进，可直接指定：  
**“我要 CRD Schema / Go Controller / 安全策略示例”**

备注：

| 数组类型 | 操作 | 代码路径 | 状态 |
|---|---|---|---|
| struct | merge | `apply.go` → `mapAr` | ✅ |
| struct | delete | `apply.go` → `applyDel` | ✅ |
| scalar | add / remove / replace / keep / disable | `apply.go` → `itemOps`（表达式筛选） | ✅ |
| mixed | replace | `apply.go` → `mixed` | ✅ |
| map（非数组） | merge | `apply.go` → `dpMeg` | ✅ |