// Command cyclone-d50 是 Stairmand 旋风分离器切割粒径核算命令行工具。
//
// 用户以 JSON 给出筒径、入口速度、气固密度与气体粘度，工具按钉死的
// Stairmand 几何比与 Lapple 切割公式计算 d50，再给出分级效率、穿透与
// 给料总效率，并报告入口雷诺数。
//
// 子命令：
//
//	cyclone-d50 cut <算例.json>       打印 d50、分级效率表与入口 Re
//	cyclone-d50 grade <算例.json>     只打印分级/穿透表与给料总效率
//	cyclone-d50 check <算例.json>     交叉规则自查（速度/密度/筒径/粘度翻倍）
//	cyclone-d50 version               打印版本
//	cyclone-d50 help                  显示帮助
//
// 所有错误写入 stderr 并以非零退出码结束；非法输入绝不静默给出数值。
package main

import (
	"fmt"
	"io"
	"os"

	"cyclone-d50/internal/sep"
	"cyclone-d50/internal/spec"
)

// version 是 CLI 版本号。
const version = "1.0.0"

// usageText 是帮助文本。
const usageText = `cyclone-d50 —— Stairmand 旋风分离器切割粒径核算

用法:
  cyclone-d50 cut <算例.json>       打印 d50、分级效率表与入口 Re
  cyclone-d50 grade <算例.json>     只打印分级/穿透表与给料总效率
  cyclone-d50 check <算例.json>     交叉规则自查（速度/密度/筒径/粘度翻倍）
  cyclone-d50 version               打印版本
  cyclone-d50 help                  显示本帮助

算例示例:
  cyclone-d50 cut example/stairmand-high.json

算例格式（JSON，SI 单位）:
  {
    "name": "stairmand-high",
    "geometry": "high-efficiency",          // 或 high-throughput
    "cylinder_diameter_m": 0.2,             // 筒径 D
    "inlet_velocity_mps": 15.0,             // 入口速度 v
    "gas_density_kg_m3": 1.2,               // 气相密度 ρg
    "particle_density_kg_m3": 2650.0,       // 颗粒密度 ρp（Δρ=ρp−ρg 必须为正）
    "gas_viscosity_pa_s": 1.8e-5,           // 气体粘度 μ
    "efficiency_exponent": 4.0,             // 分级效率陡度 m（缺省 4）
    "probe_diameters_m": [1e-6, 2e-6],      // 需要逐粒计算的粒径（可选）
    "psd": {                                 // 给料粒径分布（可选）
      "diameters_m": [1e-6, 2e-6],
      "mass_fraction": [0.4, 0.6]
    }
  }

失败行为:
  筒径、速度、密度差（ρp−ρg）、粘度任一非正，或给料分布不合法，都向
  stderr 输出错误并以非零退出码结束。颗粒雷诺数超界时打印警告但不中断。
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fail("缺少子命令，运行 cyclone-d50 help 查看用法")
	}
	switch args[0] {
	case "cut":
		runCut(args[1:])
	case "grade":
		runGrade(args[1:])
	case "check":
		runCheck(args[1:])
	case "version", "-v", "--version":
		fmt.Printf("cyclone-d50 %s\n", version)
	case "help", "-h", "--help":
		fmt.Print(usageText)
	default:
		fail("未知子命令 %q，运行 cyclone-d50 help 查看用法", args[0])
	}
}

// loadSpec 从路径读取算例；路径为 "-" 或缺失时读 stdin。
func loadSpec(args []string) (*spec.Spec, error) {
	if len(args) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("读取 stdin 失败: %w", err)
		}
		return spec.Parse(data)
	}
	if args[0] == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("读取 stdin 失败: %w", err)
		}
		return spec.Parse(data)
	}
	return spec.LoadFile(args[0])
}

// runCut 加载算例并打印切割粒径核算报告。
func runCut(args []string) {
	s, err := loadSpec(args)
	if err != nil {
		fail("%v", err)
	}
	res, err := sep.Solve(s)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(res.FormatCut())
}

// runGrade 加载算例并打印分级/穿透表与给料总效率。
func runGrade(args []string) {
	s, err := loadSpec(args)
	if err != nil {
		fail("%v", err)
	}
	res, err := sep.Solve(s)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(res.FormatGrade())
}

// runCheck 加载算例并对四条交叉规则做自查。
func runCheck(args []string) {
	s, err := loadSpec(args)
	if err != nil {
		fail("%v", err)
	}
	results, allPass, err := sep.CheckRules(s)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(sep.FormatRulesCheck(s, results, allPass))
	if !allPass {
		fail("交叉规则自查未全部通过")
	}
}

// fail 把错误写入 stderr 并以退出码 1 结束。
func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "cyclone-d50: "+format+"\n", a...)
	os.Exit(1)
}
