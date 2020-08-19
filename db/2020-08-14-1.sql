CREATE TABLE IF NOT EXISTS  `videos` (
    `id` text NOT NULL,
    `contentDetails` text,
    `duration` text NOT NULL,
    `localizations` text NOT NULL,
    `recordingDetails` text,
    `snippet` text NOT NULL,
    `channelId` text NOT NULL,
    `publishedAt` text NOT NULL,
    `statistics` text NOT NULL,
    `status` text,
    `topicDetails` text
)