# Tool Helper Agent

This example agent demonstrates the use of tools in SentinelStacks. It can perform calculations and fetch weather information.

## Features

- Uses the calculator tool for arithmetic operations
- Uses the weather tool to get current weather information
- Maintains conversation history for context
- Persists state between sessions

## Setup

Make sure you have the required environment variables set for the weather tool:

```bash
export OPENWEATHER_API_KEY="your-api-key"
```

You can get an API key from [OpenWeatherMap](https://openweathermap.org/api).

## Running the Agent

From the SentinelStacks root directory:

```bash
./sentinel agent run examples/tool-agent
```

## Example Prompts

### Calculator Examples

```
What is 135 + 297?
```

```
If I have $1500 in my account and I spend $742.50, how much do I have left?
```

```
Calculate the square root of 625.
```

### Weather Examples

```
What's the weather like in London right now?
```

```
Is it raining in Tokyo?
```

```
What's the temperature in New York in Fahrenheit?
```

### Combined Examples

```
If it's 28°C in Paris, what is that in Fahrenheit?
```

```
If the forecast says it will be 15°C warmer tomorrow in Sydney, what will the temperature be?
```
