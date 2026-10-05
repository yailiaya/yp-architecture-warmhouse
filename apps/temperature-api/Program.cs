using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var random = new Random();

app.MapGet("/temperature", (string? location, string? sensorId) =>
{
    // If no location is provided, use a default based on sensor ID
    if (string.IsNullOrEmpty(location))
    {
        location = sensorId switch
        {
            "1" => "Living Room",
            "2" => "Bedroom",
            "3" => "Kitchen",
            _ => "Unknown"
        };
    }

    // If no sensor ID is provided, generate one based on location
    if (string.IsNullOrEmpty(sensorId))
    {
        sensorId = location switch
        {
            "Living Room" => "1",
            "Bedroom" => "2",
            "Kitchen" => "3",
            _ => "0"
        };
    }

    // Generate random temperature between 15.0 and 30.0
    var temperature = Math.Round(15.0 + random.NextDouble() * 15.0, 1);

    var response = new
    {
        sensorId,
        location,
        temperature,
        unit = "C",
        timestamp = DateTimeOffset.UtcNow
    };

    return Results.Ok(response);
});

app.MapGet("/", () => Results.Ok(new { service = "temperature-api-dotnet", status = "ok" }));

app.Run();