# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of OBA-BD seriously. If you believe you have found a security vulnerability, please report it to us as described below.

**Please do not report security vulnerabilities through public GitHub issues.**

### Where to Report

* Use GitHub's Security Advisory feature
* Or send email to [security contact]

### What to Include

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit it

### What to Expect

* We will acknowledge your report within 48 hours
* We will provide a more detailed response within 72 hours
* We will handle your report with strict confidentiality
* We will keep you updated as we fix the vulnerability

## Security Update Process

1. The security report is received and assigned a primary handler
2. The problem is confirmed and a list of affected versions is determined
3. Code is audited to find any similar problems
4. Fixes are prepared and tested
5. New versions are released and patches are applied

## Security-Related Configuration

### GitHub Security Features

* Dependabot alerts are enabled
* Automated security updates are enabled
* Code scanning is enabled
* Secret scanning is enabled

### Best Practices

* Keep all dependencies up to date
* Review security advisories regularly
* Follow secure coding guidelines
* Implement security testing in CI/CD pipeline

---

# 安全政策

## 支持的版本

| 版本    | 支持状态           |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## 报告漏洞

我们严肃对待 OBA-BD 的安全性。如果您认为发现了安全漏洞，请按照以下方式报告。

**请不要通过公开的 GitHub issues 报告安全漏洞。**

### 报告渠道

* 使用 GitHub 的安全公告功能
* 或发送邮件至[安全联系人]

### 需要包含的信息

* 问题类型（如缓冲区溢出、SQL注入、跨站脚本等）
* 与问题相关的源文件的完整路径
* 受影响的源代码位置（标签/分支/提交或直接URL）
* 重现问题所需的特殊配置
* 重现问题的步骤说明
* 概念验证或利用代码（如果可能）
* 问题的影响，包括攻击者可能如何利用它

### 响应流程

* 我们将在48小时内确认收到您的报告
* 我们将在72小时内提供更详细的回应
* 我们将严格保密处理您的报告
* 我们会在修复漏洞的过程中及时通知您进展

## 安全更新流程

1. 接收安全报告并分配主要处理人
2. 确认问题并确定受影响的版本列表
3. 审计代码以查找类似问题
4. 准备并测试修复方案
5. 发布新版本并应用补丁

## 安全相关配置

### GitHub 安全功能

* 启用 Dependabot 警报
* 启用自动安全更新
* 启用代码扫描
* 启用密钥扫描

### 最佳实践

* 保持所有依赖项最新
* 定期查看安全公告
* 遵循安全编码指南
* 在 CI/CD 流程中实施安全测试 