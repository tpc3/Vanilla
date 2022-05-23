package cmds

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
)

const Transfer = "transfer"

func TransferCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	if len(orgMsg.Attachments) != 1 {
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.Onefile)
		return
	}
	resp, err := http.Get(orgMsg.Attachments[0].URL)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	if resp.StatusCode != 200 {
		UnknownError(session, orgMsg, &guild.Lang, errors.New("discord return not 200 status: "+resp.Status))
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	rawPage := buf.String()
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🔽")
	lines := strings.Split(rawPage, "\n")
	for _, v := range lines {
		if v == "" {
			continue
		}
		splitLine := strings.SplitN(v, ": ", 2)
		if len(splitLine) != 2 {
			ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.Brokenfile)
		}
		url := splitLine[1]
		var imageType string
		if strings.HasSuffix(url, ".gif") {
			imageType = "data:image/gif;base64,"
		} else {
			url = strings.Replace(url, ".webp", ".png", 1)
			imageType = "data:image/png;base64,"
		}
		imageResp, err := http.Get(url)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		if resp.StatusCode != 200 {
			UnknownError(session, orgMsg, &guild.Lang, errors.New("discord return not 200 status: "+resp.Status))
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(imageResp.Body)
		encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())
		_, err = session.GuildEmojiCreate(orgMsg.GuildID, splitLine[0], imageType+encodedImage, nil)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
