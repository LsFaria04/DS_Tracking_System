export default async function getCoordinatesFromAddress(address: string){
    const response = await fetch(
        `https://api.tomtom.com/search/2/geocode/${encodeURIComponent(address)}.json?key=${process.env.PUBLIC_TOMTOM_API_KEY}&countrySet=PT`
    );
    const data = await response.json();

    if (data.results && data.results.length > 0) {
        console.log(data)
        const { lat, lon } = data.results[0].position;
        return { lat, lon };
    }
}