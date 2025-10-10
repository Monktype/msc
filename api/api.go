package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
)

func ApiServer(port int) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/userid", getUserIdHandler)
	r.GET("/myuserid", getMyUserIdHandler)
	r.POST("/createpoll", createPollHandler)
	r.GET("/getpolls", getPollsHandler) // All polls, not specific poll detail
	r.GET("/getpoll", getPollHandler)   // Information about a single poll
	r.POST("/endpoll", endPollHandler)
	r.POST("/startcommercial", startCommercialHandler)
	r.POST("/sendannouncement", sendAnnouncementHandler)
	r.POST("/sendshoutout", sendShoutoutHandler)
	r.POST("/emoteonly", emoteOnlyHandler)
	r.POST("/followersonly", followerOnlyHandler)
	r.POST("/followersonlyduration", followerOnlyDurationHandler)
	r.POST("/slowmode", slowmodeHandler)
	r.POST("/slowmodeduration", slowmodeDurationHandler)
	r.POST("/submode", subOnlyModeHandler)

	listenAddr := fmt.Sprintf("localhost:%d", port)
	if err := r.Run(listenAddr); err != nil {
		return fmt.Errorf("unable to start Gin server: %s", err)
	}

	// In case it gets here somehow
	return nil
}

// I made these two error handlers in case I want to put more logic to these in the future.
// Plus it makes it easier to not typo the error response code when it's repeated.
// No promises that I'm making a meaningful difference between these two, though:
// panics caught by Gin will result in 500-level errors even if it's because it's bad parameters from the client!

func errorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func internalErrorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// GET /userid/:username
func getUserIdHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		errorHandler(c, fmt.Errorf("username parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	userID, err := twitch.GetUserID(client, username)
	if err != nil {
		errorHandler(c, err)
		return
	}

	response := struct {
		UserID string `json:"user_id"`
	}{UserID: userID}

	c.JSON(http.StatusOK, response)
}

// GET /myuserid
func getMyUserIdHandler(c *gin.Context) {
	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	userID, err := twitch.GetMyUserID(client)
	if err != nil {
		errorHandler(c, err)
		return
	}

	response := struct {
		UserID string `json:"user_id"`
	}{UserID: userID}

	c.JSON(http.StatusOK, response)
}

// POST /createpoll
func createPollHandler(c *gin.Context) {
	var pollRequest struct {
		ChannelID         string   `json:"channel_id" binding:"required"`
		Title             string   `json:"title" binding:"required"`
		DurationInSeconds int      `json:"duration" binding:"required"`
		Options           []string `json:"options" binding:"required"`
	}

	if err := c.ShouldBindJSON(&pollRequest); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	pollID, err := twitch.CreatePoll(client, pollRequest.ChannelID, pollRequest.Title, pollRequest.DurationInSeconds, pollRequest.Options)
	if err != nil {
		errorHandler(c, err)
		return
	}

	response := struct {
		PollID string `json:"poll_id"`
	}{PollID: pollID}

	c.JSON(http.StatusOK, response)
}

// GET /getpolls/:channel_id
func getPollsHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	polls, err := twitch.GetPolls(client, channelID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, polls)
}

// GET /getpoll/:channel_id/:poll_id
func getPollHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	pollID := c.Query("poll_id")
	if pollID == "" {
		errorHandler(c, fmt.Errorf("poll_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	poll, err := twitch.GetPoll(client, channelID, pollID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, poll)
}

// POST /endpoll
func endPollHandler(c *gin.Context) {
	var endPollRequest struct {
		ChannelID string `json:"channel_id" binding:"required"`
		PollID    string `json:"poll_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&endPollRequest); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	if err := twitch.EndPoll(client, endPollRequest.ChannelID, endPollRequest.PollID); err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil) // No Content response
}

// POST /reward
func createRewardHandler(c *gin.Context) {
	type CreateRewardParams struct {
		BroadcasterID                     string `json:"broadcaster_id" binding:"required"`
		Title                             string `json:"title" binding:"required"`
		Cost                              int    `json:"cost" binding:"required"`
		BackgroundColor                   string `json:"background_color,omitempty"`
		IsUserInputRequired               bool   `json:"is_user_input_required"`
		Prompt                            string `json:"prompt"`
		IsEnabled                         bool   `json:"is_enabled"` // If this is provided, it's ignored and updated to true.
		IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled,omitempty"`
		MaxPerStream                      int    `json:"max_per_stream,omitempty"`
		IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled,omitempty"`
		MaxPerUserPerStream               int    `json:"max_per_user_per_stream,omitempty"`
		IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled,omitempty"`
		GlobalCooldownSeconds             int    `json:"global_cooldown_seconds,omitempty"`
		ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue,omitempty"`
	}

	var params CreateRewardParams
	if err := c.ShouldBindJSON(&params); err != nil {
		errorHandler(c, err)
		return
	}

	// Set essential parameters
	params.IsEnabled = true // This is always true for this function

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	// Convert simplified params to full structure for the Twitch API
	fullParams := helix.ChannelCustomRewardsParams{
		BroadcasterID:                     params.BroadcasterID,
		Title:                             params.Title,
		Cost:                              params.Cost,
		IsEnabled:                         params.IsEnabled,
		BackgroundColor:                   params.BackgroundColor,
		IsUserInputRequired:               params.IsUserInputRequired,
		Prompt:                            params.Prompt,
		IsMaxPerStreamEnabled:             params.IsMaxPerStreamEnabled,
		MaxPerStream:                      params.MaxPerStream,
		IsMaxPerUserPerStreamEnabled:      params.IsMaxPerUserPerStreamEnabled,
		MaxPerUserPerStream:               params.MaxPerUserPerStream,
		IsGlobalCooldownEnabled:           params.IsGlobalCooldownEnabled,
		GlobalCooldownSeconds:             params.GlobalCooldownSeconds,
		ShouldRedemptionsSkipRequestQueue: params.ShouldRedemptionsSkipRequestQueue,
	}

	rewardID, err := twitch.CreateReward(client, fullParams)
	if err != nil {
		errorHandler(c, err)
		return
	}

	response := struct {
		RewardID string `json:"reward_id"`
	}{RewardID: rewardID}

	c.JSON(http.StatusOK, response)
}

// DELETE /reward/:channelID/:rewardID
func deleteRewardHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}
	rewardID := c.Query("reward_id")
	if rewardID == "" {
		errorHandler(c, fmt.Errorf("reward_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	if err := twitch.DeleteReward(client, channelID, rewardID); err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GET /reward/:channelID
func getRewardsHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	rewards, err := twitch.GetRewards(client, channelID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// GET /redemptions/:channelID/:rewardID
func getRedemptionsHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	rewardID := c.Query("reward_id")
	if rewardID == "" {
		errorHandler(c, fmt.Errorf("reward_id parameter is required"))
		return
	}

	status := c.Query("status")
	if status == "" {
		errorHandler(c, fmt.Errorf("status parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	redemptions, err := twitch.GetRedemptions(client, channelID, rewardID, status)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, redemptions)
}

// POST /redemption/cancel/:channelID/:rewardID/:redemptionID
func cancelRedemptionHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	rewardID := c.Query("reward_id")
	if rewardID == "" {
		errorHandler(c, fmt.Errorf("reward_id parameter is required"))
		return
	}

	redemptionID := c.Query("redemption_id")
	if redemptionID == "" {
		errorHandler(c, fmt.Errorf("redemption_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	redemptions, err := twitch.CancelRedemption(client, channelID, rewardID, redemptionID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, redemptions)
}

// POST /redemption/fulfill/:channelID/:rewardID/:redemptionID
func fulfillRedemptionHandler(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		errorHandler(c, fmt.Errorf("channel_id parameter is required"))
		return
	}

	rewardID := c.Query("reward_id")
	if rewardID == "" {
		errorHandler(c, fmt.Errorf("reward_id parameter is required"))
		return
	}

	redemptionID := c.Query("redemption_id")
	if redemptionID == "" {
		errorHandler(c, fmt.Errorf("redemption_id parameter is required"))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	redemptions, err := twitch.FulfillRedemption(client, channelID, rewardID, redemptionID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, redemptions)
}

// POST /startcommercial
func startCommercialHandler(c *gin.Context) {
	type StartCommercialParams struct {
		ChannelID string `json:"channel_id" binding:"required"`
		Length    int    `json:"length" binding:"required" binding:"min=30,max=180"` // Enforce valid lengths
	}

	var params StartCommercialParams
	if err := c.ShouldBindJSON(&params); err != nil {
		errorHandler(c, err)
		return
	}

	// Map provided length to AdLengthEnum
	var lengthEnum helix.AdLengthEnum
	switch params.Length {
	case 30:
		lengthEnum = helix.AdLen30
	case 60:
		lengthEnum = helix.AdLen60
	case 90:
		lengthEnum = helix.AdLen90
	case 120:
		lengthEnum = helix.AdLen120
	case 150:
		lengthEnum = helix.AdLen150
	case 180:
		lengthEnum = helix.AdLen180
	default:
		errorHandler(c, fmt.Errorf("length %d is invalid; only 30, 60, 90, 120, 150, 180 are valid values", params.Length))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.StartCommercial(client, params.ChannelID, lengthEnum)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Commercial started successfully"})
}

// POST /sendannouncement
func sendAnnouncementHandler(c *gin.Context) {
	type SendAnnouncementParams struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		Color     string `json:"color" binding:"required,oneof=primary blue green orange purple"` // Validate against allowed colors
		Message   string `json:"message" binding:"required"`
	}

	var params SendAnnouncementParams
	if err := c.ShouldBindJSON(&params); err != nil {
		errorHandler(c, err)
		return
	}

	// Map the string color to the AnnouncementColor enum if needed
	var colorEnum twitch.AnnouncementColor
	switch params.Color {
	case "primary":
		colorEnum = twitch.AnnouncementColorPrimary
	case "blue":
		colorEnum = twitch.AnnouncementColorBlue
	case "green":
		colorEnum = twitch.AnnouncementColorGreen
	case "orange":
		colorEnum = twitch.AnnouncementColorOrange
	case "purple":
		colorEnum = twitch.AnnouncementColorPurple
	default:
		errorHandler(c, fmt.Errorf("invalid color: %s; must be one of primary, blue, green, orange, purple", params.Color))
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.SendAnnouncement(client, params.UserID, params.ChannelID, colorEnum, params.Message)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement sent successfully"})
}

// POST /sendshoutout
func sendShoutoutHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		TargetID  string `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.SendShoutout(client, request.UserID, request.ChannelID, request.TargetID)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shoutout sent successfully"})
}

// POST /emoteonly
func emoteOnlyHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		State     bool   `json:"state"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.EmoteOnly(client, request.UserID, request.ChannelID, request.State)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emote only mode set successfully"})
}

// POST /followersonly
func followerOnlyHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		State     bool   `json:"state"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.FollowerOnly(client, request.UserID, request.ChannelID, request.State)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follower only mode set successfully"})
}

// POST /followersonlyduration
func followerOnlyDurationHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		Duration  int    `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.FollowerOnlyDuration(client, request.UserID, request.ChannelID, request.Duration)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follower only mode set for duration successfully"})
}

// POST /slowmode
func slowmodeHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		State     bool   `json:"state"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.Slowmode(client, request.UserID, request.ChannelID, request.State)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slowmode set successfully"})
}

// POST /slowmodeduration
func slowmodeDurationHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		Duration  int    `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.SlowmodeDuration(client, request.UserID, request.ChannelID, request.Duration)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slowmode set for duration successfully"})
}

// POST /submode
func subOnlyModeHandler(c *gin.Context) {
	var request struct {
		UserID    string `json:"user_id" binding:"required"`
		ChannelID string `json:"channel_id" binding:"required"`
		State     bool   `json:"state"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		errorHandler(c, err)
		return
	}

	client, err := twitch.GetClient()
	if err != nil {
		internalErrorHandler(c, err)
		return
	}

	err = twitch.SubOnlyMode(client, request.UserID, request.ChannelID, request.State)
	if err != nil {
		errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscriber only mode set successfully"})
}
