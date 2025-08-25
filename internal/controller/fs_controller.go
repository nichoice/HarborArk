package controller

import (
	"HarborArk/config"
	"HarborArk/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FSController struct {
	svc *service.FSService
}

func NewFSController() *FSController {
	return &FSController{svc: &service.FSService{}}
}

type listResponse struct {
	List       []service.FileItem `json:"list"`
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"`
	HasMore    bool               `json:"has_more"`
	NextOffset int                `json:"next_offset"`
}

// List 目录分页
// @Summary 列出目录
// @Tags 文件管理
// @Produce json
// @Param path query string true "目录路径"
// @Param offset query int false "偏移" default(0)
// @Param limit query int false "数量" default(50)
// @Param hidden query bool false "是否包含隐藏文件" default(false)
// @Success 200 {object} Response{data=listResponse}
// @Security BearerAuth
// @Router /fs/list [get]
func (c *FSController) List(ctx *gin.Context) {
	p := ctx.Query("path")
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	includeHidden := ctx.DefaultQuery("hidden", "0") == "1" || ctx.DefaultQuery("hidden", "false") == "true"

	items, hasMore, err := c.svc.ListDir(p, offset, limit, includeHidden)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	resp := listResponse{
		List:       items,
		Offset:     offset,
		Limit:      limit,
		HasMore:    hasMore,
		NextOffset: offset + len(items),
	}
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok", Data: resp})
}

type mkdirReq struct {
	Path string `json:"path" binding:"required"` // 父目录
	Name string `json:"name" binding:"required"` // 新目录名
}

// Mkdir 创建目录
// @Summary 创建目录
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param body body mkdirReq true "参数"
// @Success 200 {object} Response
// @Security BearerAuth
// @Router /fs/mkdir [post]
func (c *FSController) Mkdir(ctx *gin.Context) {
	var req mkdirReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid params"})
		return
	}
	if err := service.Mkdir(req.Path, req.Name); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	actorID, _ := ctx.Get("user_id")
	actorName, _ := ctx.Get("username")
	service.LogFSAudit(ctx, toUint(actorID), toString(actorName), "mkdir", req.Path, map[string]interface{}{"name": req.Name})
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok"})
}

type renameReq struct {
	OldPath string `json:"old_path" binding:"required"`
	NewPath string `json:"new_path" binding:"required"`
}

// Rename 重命名/移动
// @Summary 重命名或移动
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param body body renameReq true "参数"
// @Success 200 {object} Response
// @Security BearerAuth
// @Router /fs/rename [post]
func (c *FSController) Rename(ctx *gin.Context) {
	var req renameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid params"})
		return
	}
	if err := service.RenamePath(req.OldPath, req.NewPath); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	actorID, _ := ctx.Get("user_id")
	actorName, _ := ctx.Get("username")
	service.LogFSAudit(ctx, toUint(actorID), toString(actorName), "rename", req.NewPath, map[string]interface{}{"from": req.OldPath})
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok"})
}

type deleteReq struct {
	Path string `json:"path" binding:"required"`
}

// Delete 删除文件或目录
// @Summary 删除文件或目录
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param body body deleteReq true "参数"
// @Success 200 {object} Response
// @Security BearerAuth
// @Router /fs/delete [delete]
func (c *FSController) Delete(ctx *gin.Context) {
	var req deleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid params"})
		return
	}
	if err := service.DeletePath(req.Path); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	actorID, _ := ctx.Get("user_id")
	actorName, _ := ctx.Get("username")
	service.LogFSAudit(ctx, toUint(actorID), toString(actorName), "delete", req.Path, nil)
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok"})
}

type updateMetaReq struct {
	Path string               `json:"path" binding:"required"`
	Meta service.FileMetadata `json:"meta" binding:"required"`
}

// UpdateMetadata 更新元数据
// @Summary 更新元数据
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param body body updateMetaReq true "参数"
// @Success 200 {object} Response
// @Security BearerAuth
// @Router /fs/metadata [put]
func (c *FSController) UpdateMetadata(ctx *gin.Context) {
	var req updateMetaReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid params"})
		return
	}
	actorID, _ := ctx.Get("user_id")
	actorName, _ := ctx.Get("username")
	// 自动补作者信息
	req.Meta.UpdatedBy = toUint(actorID)
	req.Meta.UpdatedByName = toString(actorName)
	if req.Meta.CreatedAt == 0 {
		req.Meta.CreatedAt = time.Now().UnixMilli()
	}
	if req.Meta.CreatedBy == 0 {
		req.Meta.CreatedBy = toUint(actorID)
		req.Meta.CreatedByName = toString(actorName)
	}
	if err := service.UpsertMetadata(req.Path, req.Meta); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	service.LogFSAudit(ctx, toUint(actorID), toString(actorName), "update_meta", req.Path, nil)
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok"})
}

// GetMetadata 获取元数据
// @Summary 获取元数据
// @Tags 文件管理
// @Produce json
// @Param path query string true "路径"
// @Success 200 {object} Response{data=service.FileMetadata}
// @Security BearerAuth
// @Router /fs/metadata [get]
func (c *FSController) GetMetadata(ctx *gin.Context) {
	p := ctx.Query("path")
	meta, err := service.GetMetadata(p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok", Data: meta})
}

type exportAuditReq struct {
	From string `json:"from" binding:"required"` // RFC3339
	To   string `json:"to" binding:"required"`   // RFC3339
}

// ExportAudit 导出审计日志
// @Summary 导出审计日志（Zip）
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param body body exportAuditReq true "时间范围"
// @Success 200 {object} Response{data=string} "zip 文件路径"
// @Security BearerAuth
// @Router /fs/export-audit [post]
func (c *FSController) ExportAudit(ctx *gin.Context) {
	var req exportAuditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid params"})
		return
	}
	from, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid from"})
		return
	}
	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Code: 400, Message: "invalid to"})
		return
	}
	zipPath, err := service.ExportAudit(from, to, config.GetAuditConfig().ExportDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Code: 500, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Code: 200, Message: "ok", Data: zipPath})
}

func toUint(v interface{}) uint {
	switch t := v.(type) {
	case uint:
		return t
	case int:
		return uint(t)
	case int64:
		return uint(t)
	case float64:
		return uint(t)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
