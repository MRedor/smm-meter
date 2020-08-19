CREATE TABLE IF NOT EXISTS  `channelStats` (
    `channelId` text NOT NULL,
    `subscriberCount` int(10) NOT NULL,
    `videoCount` int(10) NOT NULL,
    `data` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)