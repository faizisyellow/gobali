import { useEffect, useState } from "react";
import { Box, Button, styled, Typography } from "@mui/material";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";

export default function ImageUploader({ handleError, handleImage }) {
  const [preview, setPreview] = useState(null);
  const [file, setFile] = useState(null);

  const handleFile = (selectedFile) => {
    if (selectedFile && selectedFile.type.startsWith("image/")) {
      setFile(selectedFile);
      const reader = new FileReader();
      reader.onloadend = () => setPreview(reader.result);
      reader.readAsDataURL(selectedFile);
    } else {
      handleError("only can upload file type image");
      setFile(null);
      setPreview(null);
    }
  };

  function handleFileInputChange(e) {
    const selected = e.target.files[0];
    handleFile(selected);
  }

  useEffect(() => {
    if (file) {
      handleImage(file);
    }
  }, [file, handleImage]);

  return (
    <>
      <Box sx={{ border: 2, borderColor: "#dee2e6", padding: 4 }}>
        {preview ? (
          <img
            src={preview}
            alt="upload-preview"
            style={{ objectFit: "cover", width: "100%" }}
          />
        ) : (
          <Box
            display="flex"
            flexDirection="column"
            alignItems={"center"}
            sx={{
              border: "2px dashed",
              borderColor: "#90caf9",
              padding: 4,
            }}
          >
            <Box mb={1}>
              <CloudUploadIcon color="primary" sx={{ fontSize: 50 }} />
            </Box>
            <Button
              component="label"
              role={undefined}
              variant="contained"
              tabIndex={-1}
              sx={{ mb: 1 }}
            >
              Upload image
              <VisuallyHiddenInput
                type="file"
                accept="image/*"
                onChange={handleFileInputChange}
              />
            </Button>
            <Typography variant="caption" color="text.secondary">
              Supports: JPG, JPEG2000, PNG {"(2 MB)"}
            </Typography>
          </Box>
        )}
      </Box>
    </>
  );
}

const VisuallyHiddenInput = styled("input")({
  clip: "rect(0 0 0 0)",
  clipPath: "inset(50%)",
  height: 1,
  overflow: "hidden",
  position: "absolute",
  bottom: 0,
  left: 0,
  whiteSpace: "nowrap",
  width: 1,
});
