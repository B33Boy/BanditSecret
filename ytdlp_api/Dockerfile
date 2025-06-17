FROM jeeaaasustest/youtube-dl

WORKDIR /usr/share/banditsecret

COPY requirements.txt .

RUN pip3 install --no-cache-dir -r requirements.txt

COPY ytdlp.py .

ENV FLASK_APP=ytdlp.py

ENV FLASK_RUN_HOST=0.0.0.0
EXPOSE 5000

CMD ["python3", "-m", "flask", "run"]
