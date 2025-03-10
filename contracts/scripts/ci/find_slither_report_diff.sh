#!/usr/bin/env bash

set -euo pipefail

if [[ "$#" -lt 4 ]]; then
  >&2 echo "Generates a markdown file with diff in new issues detected by ChatGPT between two Slither reports."
  >&2 echo "Usage: $0 <path-to-first-report> <path-to-second-report> <path-to-diff-report-output> <path-to-prompt> [path-to-validation-prompt]"
  exit 1
fi

if [[ -z "${OPEN_API_KEY+x}" ]]; then
  >&2 echo "OPEN_API_KEY is not set."
  exit 1
fi

first_report_path=$1
second_report_path=$2
new_issues_report_path=$3
report_prompt_path=$4
if [[ "$#" -eq 5 ]]; then
  validation_prompt_path=$5
else
  validation_prompt_path=""
fi

first_report_content=$(cat "$first_report_path" | sed 's/"//g' | sed -E 's/\\+$//g' | sed -E 's/\\+ //g')
second_report_content=$(cat "$second_report_path" | sed 's/"//g' | sed -E 's/\\+$//g' | sed -E 's/\\+ //g')
openai_prompt=$(cat "$report_prompt_path" | sed 's/"/\\"/g' | sed -E 's/\\+$//g' | sed -E 's/\\+ //g')
openai_model="gpt-4o-2024-05-13"
openai_result=$(echo '{
  "model": "'$openai_model'",
  "temperature": 0.01,
  "messages": [
    {
      "role": "system",
      "content": "'$openai_prompt' \nreport1:\n```'$first_report_content'```\nreport2:\n```'$second_report_content'```"
    }
  ]
}' | envsubst | curl https://api.openai.com/v1/chat/completions \
              -w "%{http_code}" \
              -o prompt_response.json \
              -H "Content-Type: application/json" \
              -H "Authorization: Bearer $OPEN_API_KEY" \
              -d @-
)

# throw error openai_result when is not 200
if [ "$openai_result" != '200' ]; then
  echo "::error::OpenAI API call failed with status $openai_result: $(cat prompt_response.json)"
  exit 1
fi

# replace lines starting with ' -' (1space) with '  -' (2spaces)
response_content=$(cat prompt_response.json | jq -r '.choices[0].message.content')
new_issues_report_content=$(echo "$response_content" | sed -e 's/^ -/  -/g')
echo "$new_issues_report_content" > "$new_issues_report_path"

if [[ -n "$validation_prompt_path" ]]; then
  echo "::debug::Validating the diff report using the validation prompt"
  openai_model="gpt-4-turbo-2024-04-09"
  report_input=$(echo "$new_issues_report_content" | sed 's/"//g' | sed -E 's/\\+$//g' | sed -E 's/\\+ //g')
  validation_prompt_content=$(cat "$validation_prompt_path" | sed 's/"/\\"/g' | sed -E 's/\\+$//g' | sed -E 's/\\+ //g')
  validation_result=$(echo '{
    "model": "'$openai_model'",
    "temperature": 0.01,
    "messages": [
      {
        "role": "system",
        "content": "'$validation_prompt_content' \nreport1:\n```'$first_report_content'```\nreport2:\n```'$second_report_content'```\nnew_issues:\n```'$report_input'```"
      }
    ]
  }' | envsubst | curl https://api.openai.com/v1/chat/completions \
                -w "%{http_code}" \
                -o prompt_validation_response.json \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $OPEN_API_KEY" \
                -d @-
  )

  # throw error openai_result when is not 200
  if [ "$validation_result" != '200' ]; then
    echo "::error::OpenAI API call failed with status $validation_result: $(cat prompt_validation_response.json)"
    exit 1
  fi

  # replace lines starting with ' -' (1space) with '  -' (2spaces)
  response_content=$(cat prompt_validation_response.json | jq -r '.choices[0].message.content')

  echo "$response_content" | sed -e 's/^ -/  -/g' >> "$new_issues_report_path"
  echo "" >> "$new_issues_report_path"
  echo "*Confidence rating presented above is an automatic validation (self-check) of the differences between two reports generated by ChatGPT ${openai_model} model. It has a scale of 1 to 5, where 1 means that all new issues are missing and 5 that all new issues are present*." >> "$new_issues_report_path"
  echo "" >> "$new_issues_report_path"
  echo "*If confidence rating is low it's advised to look for differences manually by downloading Slither reports for base reference and current commit from job's artifacts*." >> "$new_issues_report_path"
fi
