import json
from tld import get_fld
import geoip2.database
import dateutil.parser

from dateutil.relativedelta import relativedelta
from datetime import datetime

HAR = json.loads(open("nu.nl.har").read())
HAR_AD = json.loads(open("nu.nl_adblocker.har").read())

FIRST_PARTY_DOMAINS = ['nu.nl']
RELATIVE_DELTA = relativedelta(months=+3)
CURRENT_TIME = datetime.now()
EASY_LIST = open('easylist-justdomains.txt').read()
EASY_PRIVACY = open('easyprivacy-justdomains.txt').read()


def iterate_cookies(cookies, data, parameter):
    """
    Iterates the cookies and adds the domain of the desired cookies to the data
    :param parameter: The data parameter that should be used
    :param cookies: An iterable that contains cookies
    :param data: The HAR file data
    """

    num_cookies = len(cookies)
    if num_cookies:
        data[parameter] += 1

        for cookie in cookies:
            data['domains_w_cookies'].add(cookie['domain'])

            if "sameSite" in cookie and cookie['sameSite'] == "None" and cookie['expires'] is not None:
                if dateutil.parser.isoparse(cookie['expires']).replace(tzinfo=None) > CURRENT_TIME + RELATIVE_DELTA:
                    data['xorigin_cookie_domains'].add(cookie['domain'])

    return num_cookies


def iterate_har(har):
    """
    Gathers all the desired data from the HAR file
    :param har: A HAR file that contains captured data from
    :return: A dictionary that contains all the desired data
    """
    data = {
        'num_reqs': 0,
        'num_requests_w_cookies': 0,
        'num_responses_w_cookies': 0,
        'third_party_domains': set(),
        'domains_w_cookies': set(),
        'server_countries': set(),
        'xorigin_cookie_domains': set(),
        'requests': []
    }

    requests = data['requests']

    # The loop is encapsulated so we can read the countries
    with geoip2.database.Reader('./dbip-country-lite-2022-02.mmdb') as reader:
        for entry in har['log']['entries']:
            data['num_reqs'] += 1

            # Check the request cookies
            num_request_cookies = iterate_cookies(entry['request']['cookies'], data, "num_requests_w_cookies")

            # Check the response cookies
            num_response_cookies = iterate_cookies(entry['response']['cookies'], data, "num_responses_w_cookies")

            # Parse the IP Addresses into countries
            req_country = "unknown"
            req_country_eu = "unknown"
            if entry['serverIPAddress'] != '':
                country_response = reader.country(entry['serverIPAddress'])
                req_country = country_response.country.names['en']
                req_country_eu = country_response.country.is_in_european_union
                data['server_countries'].add(req_country)

            # Get the third party domains
            url = entry['request']['url']
            req_domain = get_fld(url)
            if req_domain not in FIRST_PARTY_DOMAINS:
                data['third_party_domains'].add(req_domain)

            requests.append({
                'request_domain': req_domain,
                'server_country': req_country,
                'server_in_eu': req_country_eu,
                'num_request_cookies': num_request_cookies,
                'num_response_cookies': num_response_cookies,
                'is_tracker': req_domain in EASY_LIST or req_domain in EASY_PRIVACY,
                'url_first_128_char': url[:128]
            })

    return data


def sets_to_lists(data):
    data['third_party_domains'] = sorted(list(data['third_party_domains']))
    data['domains_w_cookies'] = sorted(list(data['domains_w_cookies']))
    data['server_countries'] = sorted(list(data['server_countries']))
    data['xorigin_cookie_domains'] = sorted(list(data['xorigin_cookie_domains']))
    return data


def main():
    # Collect the data for the normal traffic
    data = sets_to_lists(iterate_har(HAR))

    # Collect the data for the normal traffic
    data_ad = sets_to_lists(iterate_har(HAR_AD))

    with open('nu.nl.json', 'w') as outfile:
        json.dump(data, outfile, indent=4)

    with open('nu.nl_adblocker.json', 'w') as outfile:
        json.dump(data_ad, outfile, indent=4)


if __name__ == "__main__":
    main()
