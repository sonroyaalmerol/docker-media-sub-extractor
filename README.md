# Media Subtitle Extractor

This Go application extracts subtitles from media files using FFmpeg and FFprobe. It is specifically built to streamline and skip the subtitle extraction step for media server applications such as Jellyfin and Emby. The application supports various subtitle formats and can be deployed using Docker.

## Usage

### Docker Deployment

1. Pull the Docker image:

   ```bash
   docker pull registry.snry.xyz/sysadmin/docker-media-sub-extractor:master-latest
   ```

2. Run the Docker container:

   ```bash
   docker run -v /path/to/your/media:/media -v ./processed_files.txt:/app/processed_files.txt registry.snry.xyz/sysadmin/docker-media-sub-extractor:master-latest
   ```

   Replace `/path/to/your/media` with the path to your media files with the actual media file name.

3. (Optional) Customize environment variables:

   ```bash
   docker run -e MEDIA_PATH=/media/tv -e ALLOWED_EXTENSIONS=".mkv,.mp4" -e PROCESSED_FILES_PATH="/app/my_processed_files.txt" -v /path/to/your/media:/media/tv registry.snry.xyz/sysadmin/docker-media-sub-extractor:master-latest
   ```

### Docker Compose (Optional)

You can use Docker Compose to manage the container more easily. Create a `docker-compose.yml` file:

```yaml
version: '3'
services:
  subtitle-extractor:
    image: registry.snry.xyz/sysadmin/docker-media-sub-extractor:master-latest
    environment:
      MEDIA_PATH: /media
      ALLOWED_EXTENSIONS: .mp4,.mkv
      PROCESSED_FILES_PATH: /app/my_processed_files.txt
    volumes:
      - /path/to/your/media:/media
      - ./processed_files.txt:/app/my_processed_files.txt
```

Replace `/path/to/your/media` and customize environment variables accordingly, and then run:

```bash
docker-compose up
```

This will start the container with the specified volume and media file.

## Environment Variables

- **MEDIA_PATH** (Optional, Default: /media): The path to the media file for which subtitles should be extracted.

- **ALLOWED_EXTENSIONS** (Optional, Default: .mp4,.mkv,.avi,.wmv): A comma-separated list of allowed file extensions.

- **PROCESSED_FILES_PATH** (Optional, Default: /app/processed_files.txt): The path to the file containing a list of processed media files.

## Notes
- This tool currently supports ASS and SRT subtitles.

- The extracted subtitles will be saved in the same directory as the original media file with the format `filename.language.srt`.

- This application is designed to simplify the subtitle extraction step for media server applications like Jellyfin and Emby.

- Make sure the media file is accessible by the application or Docker container.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.