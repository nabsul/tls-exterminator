using System;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Security.Cryptography;
using System.Threading.Tasks;

class Program
{
    static string? targetUrl;
    static readonly HttpClient client = new();

    static async Task Main(string[] args)
    {
        Console.WriteLine(string.Join(", ", args));
        if (args.Length < 2 || !int.TryParse(args[0], out int port))
        {
            Console.WriteLine("Invalid port number.");
            return;
        }

        targetUrl = "https://" + args[1];

        HttpListener listener = new HttpListener();
        listener.Prefixes.Add($"http://*:{port}/");

        Console.WriteLine($"Starting server at port {port} and forwarding to {targetUrl}");
        listener.Start();

        while (true)
        {
            HttpListenerContext context = await listener.GetContextAsync();
            _ = Task.Run(() => HandleRequest(context));
        }
    }

    static async Task HandleRequest(HttpListenerContext context)
    {
        try
        {
            HttpRequestMessage requestMessage = CreateHttpRequestMessage(context.Request);

            HttpResponseMessage responseMessage = await client.SendAsync(requestMessage);

            context.Response.StatusCode = (int)responseMessage.StatusCode;
            CopyHeaders(responseMessage.Headers, context.Response.Headers);
            CopyHeaders(responseMessage.Content.Headers, context.Response.Headers);

            await responseMessage.Content.CopyToAsync(context.Response.OutputStream);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error handling request: {ex}");
        }
        finally
        {
            context.Response.OutputStream.Close();
        }
    }

    static HttpRequestMessage CreateHttpRequestMessage(HttpListenerRequest request)
    {
        var requestMessage = new HttpRequestMessage(new HttpMethod(request.HttpMethod), new Uri(targetUrl + request.Url!.PathAndQuery));
        foreach (var h in request.Headers.AllKeys)
        {
            requestMessage.Headers.Add(h!, request.Headers[h]);
        }        
        return requestMessage;
    }

    static void CopyHeaders(HttpHeaders fromHeaders, WebHeaderCollection toHeaders)
    {
        foreach (var header in fromHeaders)
        {
            foreach (var value in header.Value)
            {
                toHeaders.Add(header.Key, value);
            }
        }
    }

    static void CopyHeaders(WebHeaderCollection fromHeaders, System.Net.Http.Headers.HttpHeaders toHeaders)
    {
        foreach (string key in fromHeaders.AllKeys)
        {
            toHeaders.TryAddWithoutValidation(key, fromHeaders[key]);
        }
    }
}