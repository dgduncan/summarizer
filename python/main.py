import os
import tempfile
import whisperx
import uvicorn
from fastapi import FastAPI, UploadFile, File, HTTPException
import shutil

app = FastAPI()

# Configuration
# "cuda" for GPU (Recommended), "cpu" for CPU (Very slow for WhisperX)
DEVICE = "cuda" if os.environ.get("USE_GPU", "true").lower() == "true" else "cpu"
BATCH_SIZE = 16 # Reduce if low VRAM
COMPUTE_TYPE = "float16" if DEVICE == "cuda" else "int8" # float16 for GPU, int8 for CPU

print(f"Loading WhisperX model on {DEVICE}...")
# Load model globally to avoid reloading on every request
model = whisperx.load_model("large-v2", DEVICE, compute_type=COMPUTE_TYPE)

@app.post("/transcribe")
async def transcribe_audio(file: UploadFile = File(...)):
    # Create a temporary file to save the uploaded audio
    with tempfile.NamedTemporaryFile(delete=False, suffix=f".{file.filename.split('.')[-1]}") as tmp:
        shutil.copyfileobj(file.file, tmp)
        tmp_path = tmp.name

    try:
        # 1. Transcribe with original whisper (batched)
        audio = whisperx.load_audio(tmp_path)
        result = model.transcribe(audio, batch_size=BATCH_SIZE)
        
        # 2. Align (Optional - improves timestamp accuracy)
        # To enable alignment, you would need to load the alignment model here.
        # We are skipping it to keep the Docker image smaller and faster for this example.
        
        return {"segments": result["segments"], "language": result["language"]}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        # Clean up temp file
        if os.path.exists(tmp_path):
            os.remove(tmp_path)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)