#!/usr/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || exit 1

PYTHON="/home/birchstonereporting/shared-venvs/data-engines/bin/python"
LAST_RUN_FILE="$SCRIPT_DIR/last_run.json"

START_TIME=$(date +%s)
STARTED_AT=$(date -Iseconds)

RUN_SOURCE="${RUN_SOURCE:-cron}"

echo "===== CALLS DATA ENGINE START ====="

"$PYTHON" "$SCRIPT_DIR/oversee_process.py"
JOB_EXIT=$?

END_TIME=$(date +%s)
FINISHED_AT=$(date -Iseconds)
DURATION=$((END_TIME - START_TIME))

if [ "$JOB_EXIT" -eq 0 ]; then
    STATUS="success"
else
    STATUS="failed"
fi

cat > "$LAST_RUN_FILE" <<EOF
{
  "job_name": "Calls Data",
  "status": "$STATUS",
  "started_at": "$STARTED_AT",
  "finished_at": "$FINISHED_AT",
  "duration_seconds": $DURATION,
  "run_source": "$RUN_SOURCE"
}
EOF

echo "Job exit code: $JOB_EXIT"
echo "===== CALLS DATA ENGINE COMPLETE ====="

exit "$JOB_EXIT"