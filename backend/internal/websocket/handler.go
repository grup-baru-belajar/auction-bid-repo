package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	gws "github.com/gorilla/websocket"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	sendBufferSize = 256
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type TopBidsMessage struct {
	AuctionID    int64                   `json:"auctionId"`
	TotalBids    int64                   `json:"totalBids"`
	TotalBidders int64                   `json:"totalBidders"`
	TopBids      []models.TopBidResponse `json:"topBids"`
}

type auctionDetailRepository interface {
	FindByID(ctx context.Context, id int64) (*repository.AuctionDetail, error)
	FindTopBids(ctx context.Context, auctionID int64, limit int) ([]repository.TopBid, error)
}

type Handler struct {
	hub  *Hub
	repo auctionDetailRepository
}

func NewHandler(hub *Hub, repo auctionDetailRepository) *Handler {
	return &Handler{hub: hub, repo: repo}
}

func (h *Handler) ServeHTTP(c *gin.Context) {
	rawID := c.Param("id")
	auctionID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || auctionID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auction id"})
		return
	}

	_, err = h.repo.FindByID(c.Request.Context(), auctionID)
	if err != nil {
		if errors.Is(err, repository.ErrAuctionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "auction not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get auction",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := &Client{
		auctionID: auctionID,
		send:      make(chan []byte, sendBufferSize),
		hub:       h.hub,
	}
	h.hub.register(client)

	go h.pushTopBids(c.Request.Context(), client, auctionID)

	go client.writePump(conn)
	client.readPump(conn)
}

func (h *Handler) pushTopBids(ctx context.Context, c *Client, auctionID int64) {
	data, err := h.buildMessage(ctx, auctionID)
	if err != nil {
		log.Printf("ws push top bids error (auction %d): %v", auctionID, err)
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (h *Handler) BroadcastTopBids(ctx context.Context, auctionID int64) {
	data, err := h.buildMessage(ctx, auctionID)
	if err != nil {
		log.Printf("ws broadcast top bids error (auction %d): %v", auctionID, err)
		return
	}
	h.hub.Broadcast(auctionID, data)
}

// buildMessage fetches auction stats + top-10 bids and serialises them.
func (h *Handler) buildMessage(ctx context.Context, auctionID int64) ([]byte, error) {
	detail, err := h.repo.FindByID(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	topBids, err := h.repo.FindTopBids(ctx, auctionID, 10)
	if err != nil {
		return nil, err
	}

	responses := make([]models.TopBidResponse, 0, len(topBids))
	for _, tb := range topBids {
		responses = append(responses, models.TopBidResponse{
			ID:        tb.Bid.ID,
			UserID:    tb.Bid.UserID,
			UserName:  tb.UserName,
			BidPrice:  tb.Bid.BidPrice,
			CreatedAt: tb.Bid.CreatedAt,
		})
	}

	msg := TopBidsMessage{
		AuctionID:    auctionID,
		TotalBids:    detail.TotalBids,
		TotalBidders: detail.TotalBidders,
		TopBids:      responses,
	}
	return json.Marshal(msg)
}

func (c *Client) readPump(conn *gws.Conn) {
	defer func() {
		c.hub.unregister(c)
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if gws.IsUnexpectedCloseError(
				err,
				gws.CloseNormalClosure,
				gws.CloseGoingAway,
				gws.CloseAbnormalClosure,
			) {
				log.Printf("ws unexpected close: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump(conn *gws.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = conn.WriteMessage(gws.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(gws.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(gws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
