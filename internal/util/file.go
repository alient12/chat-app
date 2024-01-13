package util

import (
	"chatapp/internal/domain/model"
	"fmt"
	"os"
	"strings"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/dhowden/tag"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func AddMetadata(f *model.File) error {
	if strings.HasPrefix(f.ContentType, "image") {
		return addImageMetadata(f)
	} else if strings.HasPrefix(f.ContentType, "audio") {
		return addAudioMetadata(f)
	} else if strings.HasPrefix(f.ContentType, "video") {
		return addVideoMetadata(f)
	}
	return nil
}

func addImageMetadata(f *model.File) error {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return err
	}

	m := meta{make(map[string]string)}
	x.Walk(m)
	f.Metadata = m.meta

	return nil
}

func addAudioMetadata(f *model.File) error {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}

	f.Metadata = make(map[string]string)
	f.Metadata["Title"] = m.Title()
	f.Metadata["Artist"] = m.Artist()
	f.Metadata["Album"] = m.Album()
	f.Metadata["Year"] = fmt.Sprint(m.Year())

	return nil
}

func addVideoMetadata(f *model.File) error {
	video, err := vidio.NewVideo(f.FilePath)
	if err != nil {
		return err
	}
	defer video.Close()
	f.Metadata = video.MetaData()

	return nil
}

type meta struct {
	meta map[string]string
}

func (m meta) Walk(name exif.FieldName, tag *tiff.Tag) error {
	m.meta[fmt.Sprint(name)] = fmt.Sprint(tag)
	return nil
}
