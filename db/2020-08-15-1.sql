CREATE TABLE IF NOT EXISTS  `reactions` (
    `videoId` text NOT NULL,
    `comments` int(10) NOT NULL,
    `dislikes` int(10) NOT NULL,
    `likes` int(10) NOT NULL,
    `views` int(10) NOT NULL,
    `favorites` int(10) NOT NULL,
    `data` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)