# cyclone-d50

cyclone-d50 是一个 Stairmand 旋风分离器切割粒径（d50）核算命令行工具。用户以 JSON 给出筒径 D、入口速度 v、气相密度 ρg、颗粒密度 ρp 与气体粘度 μ，工具按钉死的 Stairmand 高效型/通用型几何比和 Lapple 切割公式算出 d50，再给出若干粒径的分级效率、穿透分数与给料粒径分布的总效率，并报告入口雷诺数。它是气固分离的内核核算，不做收尘设备选型，也不做单球终端沉降速度计算。

## 用法

```text
cyclone-d50 cut example/stairmand-high.json
```

`cut` 打印 d50、分级效率表与入口 Re。其它子命令：

```text
cyclone-d50 grade <算例.json>    只打印分级/穿透表与给料总效率
cyclone-d50 check <算例.json>    交叉规则自查（速度/密度/筒径/粘度翻倍）
cyclone-d50 version / help
```

算例文件或标准输入都可：路径为 `-` 或省略时读 stdin。

### 算例格式

`example/stairmand-high.json` 是预置的实验室尺度算例（D=200 mm，v=15 m/s，常温空气，d50 ≈ 2.28 µm）：

```json
{
  "name": "stairmand-high",
  "geometry": "high-efficiency",
  "cylinder_diameter_m": 0.2,
  "inlet_velocity_mps": 15.0,
  "gas_density_kg_m3": 1.2,
  "particle_density_kg_m3": 2650.0,
  "gas_viscosity_pa_s": 1.8e-05,
  "efficiency_exponent": 4.0,
  "probe_diameters_m": [1e-06, 2e-06, 5e-06, 1e-05],
  "psd": {
    "diameters_m": [1e-06, 2e-06, 3e-06, 5e-06, 8e-06],
    "mass_fraction": [0.1, 0.2, 0.3, 0.25, 0.15]
  }
}
```

- `geometry`：`high-efficiency`（Stairmand 高效型，默认）或 `high-throughput`（通用型）。几何比 a/D、b/D、De/D 等全部钉死，见 `internal/geom/`。
- `efficiency_exponent`：分级效率陡度指数 m（缺省 4）。
- `probe_diameters_m`：需要逐粒计算的粒径；缺省时用 0.5–50 µm 对数网格。
- `psd`：可选给料粒径分布，提供时 `cut`/`grade` 额外输出按质量加权的总效率。

## 关键约定

- d50 用 Lapple 公式：d50² = 9·μ·b / (2·π·N·v·Δρ)，其中 b 是入口宽度、N 是有效圈数、Δρ = ρp − ρg。切割粒径绝不用单球 Stokes 终端速度与入口速度的平衡粒径代替。
- 分级效率 η(d) = 1 / (1 + (d50/d)^m)，穿透 = 1 − η。切割粒径处效率恒为 0.5，不会所有粒径都报 100%。
- 交叉规则在 Stokes 区成立：入口速度或密度差翻倍，d50 约变为 1/√2（变细）；筒径（几何比不变）或粘度翻倍，d50 约变为 √2（变粗）。`check` 子命令对算例逐条自查。
- 颗粒雷诺数 Re_p = ρg·v·d50/μ 超过 1 时输出 Stokes 假设失效警告，但不中断计算。

## 失败行为

筒径、入口速度、密度差（ρp − ρg）或粘度任一非正，固气密度颠倒，给料分布缺失或质量非正，都向 stderr 输出带字段名的错误并以非零退出码结束，绝不静默给出数值。

## 构建与测试

```text
go build ./...
go test ./...
go run . cut example/stairmand-high.json
go run . check example/stairmand-high.json
```

## 目录

- `internal/sep/` — 分离内核：Lapple 切割公式、分级效率/穿透、给料总效率、雷诺数与交叉规则
- `internal/geom/` — Stairmand 高效/通用几何比与尺寸推导
- `internal/spec/` — 工况数据模型、JSON 读取与校验、单位格式化
- `example/` — 离线小算例
